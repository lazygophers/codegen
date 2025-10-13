package i18n

import (
	"errors"
	"fmt"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"github.com/pterm/pterm"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type TranslateConfig struct {
	SrcFile            string
	SrcLang            *Language
	Langs              []*Language
	OverwriteKeyPrefix []string
	Overwrite          bool

	AutoTran *state.CfgI18nAutoTran

	Localizer  Localizer
	Translator Translator
}

func Translate(c *TranslateConfig) error {
	srcLocalize, err := LoadLocalize(c.SrcFile, c.Localizer)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = updateAutoTran(c.AutoTran, filepath.Dir(c.SrcFile))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var cache *TranCache
	if c.AutoTran.EnableRecord {
		cache, err = NewTranCache(c.AutoTran.RecordPath)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		defer cache.Close()
	}

	for _, dstLang := range c.Langs {
		if dstLang.Lang == c.SrcLang.Lang {
			continue
		}

		pterm.Info.Printfln("try translate to %s", dstLang)

		log.Infof("load from %s", filepath.Join(filepath.Dir(c.SrcFile), dstLang.Lang+filepath.Ext(c.SrcFile)))

		dstLocalize, err := LoadLocalize(filepath.Join(filepath.Dir(c.SrcFile), dstLang.Lang+filepath.Ext(c.SrcFile)), c.Localizer)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		var tx *TranTx
		if cache != nil {
			tx, err = cache.Begin()
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			// Setup deferred rollback on error
			defer func() {
				if err != nil && tx != nil {
					_ = tx.Rollback()
				}
			}()
		}

		err = translate("", srcLocalize, dstLang, dstLocalize, c, tx)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = SaveLocalize(filepath.Join(filepath.Dir(c.SrcFile), dstLang.Lang+filepath.Ext(c.SrcFile)), dstLocalize, c.Localizer)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		if tx != nil {
			err = tx.Commit()
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
		}
	}

	return nil
}

type (
	BeforeTranslate func(lang *Language, key, source, dest string) (skip bool)
	AfterTranslate  func(lang *Language, key, source, dest string)
)

var (
	beforeTranslate BeforeTranslate
	afterTranslate  AfterTranslate
)

func RegisterBeforeTranslate(fn BeforeTranslate) {
	beforeTranslate = fn
}

func RegisterAfterTranslate(fn AfterTranslate) {
	afterTranslate = fn
}

// convertToStringAnyMap converts various map types to map[string]any using reflection
func convertToStringAnyMap(v any) (map[string]any, error) {
	// Fast path for already correct type
	if m, ok := v.(map[string]any); ok {
		return m, nil
	}

	// Use reflection for generic handling
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Map {
		return nil, errors.New(fmt.Sprintf("expected map type, got %T", v))
	}

	sub := make(map[string]any)
	iter := val.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()

		// Convert key to string
		var keyStr string
		switch k.Kind() {
		case reflect.String:
			keyStr = k.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			keyStr = strconv.FormatInt(k.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			keyStr = strconv.FormatUint(k.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			keyStr = strconv.FormatFloat(k.Float(), 'f', -1, 64)
		default:
			keyStr = candy.ToString(k.Interface())
		}

		sub[keyStr] = v.Interface()
	}

	return sub, nil
}

// shouldTranslate checks if a key needs to be translated
func shouldTranslate(c *TranslateConfig, dstLang *Language, parent, key string, srcValue any, dstValue any, tx *TranTx) bool {
	fullKey := parent + "." + key

	// Check beforeTranslate hook
	if beforeTranslate != nil {
		var dstStr string
		if dstValue != nil {
			dstStr = candy.ToString(dstValue)
		}
		if !beforeTranslate(dstLang, fullKey, candy.ToString(srcValue), dstStr) {
			return true
		}
	}

	// Check overwrite mode
	if c.Overwrite {
		log.Warnf("need tran %s because of overwrite", key)
		return true
	}

	// Check if key doesn't exist
	if dstValue == nil {
		log.Warnf("need tran %s because of not exist", key)
		return true
	}

	// Check overwrite key prefix
	for _, prefix := range c.OverwriteKeyPrefix {
		if strings.HasPrefix(fullKey, prefix) {
			log.Warnf("need tran %s because of overwrite key prefix %s", key, prefix)
			return true
		}
	}

	// Check source text changed via cache
	if tx != nil {
		ok, err := tx.NeedTran(dstLang.Lang, fullKey, srcValue)
		if err != nil {
			log.Errorf("err:%v", err)
		} else if ok {
			log.Warnf("need tran %s because of source text changed", key)
			return true
		}
	}

	return false
}

func translate(parent string, srcLocalize map[string]any, dstLang *Language, dstLocalize map[string]any, c *TranslateConfig, tx *TranTx) (err error) {
	toDstSub := func(k string) (map[string]any, error) {
		var sub map[string]any
		var err error
		if v, ok := dstLocalize[k]; ok {
			sub, err = convertToStringAnyMap(v)
			if err != nil {
				return nil, err
			}
		} else {
			sub = make(map[string]any)
		}

		dstLocalize[k] = sub
		return sub, nil
	}

	for k, v := range srcLocalize {
		log.Infof("try translate %s from %s to %s", parent+"."+k, c.SrcLang, dstLang)

		switch x := v.(type) {
		case string:
			// Skip empty or whitespace-only strings
			if strings.TrimSpace(x) == "" {
				log.Infof("skip empty string for key %s", parent+"."+k)
				if _, exists := dstLocalize[k]; !exists {
					dstLocalize[k] = x
				}
				continue
			}

			if shouldTranslate(c, dstLang, parent, k, v, dstLocalize[k], tx) {
				log.Infof("try translate [%s]%s from %s to %s", parent+"."+k, x, c.SrcLang, dstLang)

				var targetList []string
				var hasError bool

				for _, s := range strings.Split(x, "\n") {
					if s == "" {
						targetList = append(targetList, s)
						continue
					}

					// Skip whitespace-only lines
					if strings.TrimSpace(s) == "" {
						targetList = append(targetList, s)
						continue
					}

					target, err := c.Translator.Translate(c.SrcLang.Tag, dstLang.Tag, s)
					if err != nil {
						log.Errorf("err:%v", err)
						hasError = true
						break
					}
					targetList = append(targetList, target)
				}

				if !hasError {
					dstLocalize[k] = strings.Join(targetList, "\n")
					if afterTranslate != nil {
						afterTranslate(dstLang, parent+"."+k, candy.ToString(v), strings.Join(targetList, "\n"))
					}

					if tx != nil {
						err := tx.Update(dstLang.Lang, parent+"."+k, srcLocalize[k])
						if err != nil {
							log.Errorf("err:%v", err)
							return err
						}
					}

					pterm.Success.Printfln("key:%s\nfrom(%s) %s\nto(%s) %s",
						parent+"."+k,
						c.SrcLang,
						srcLocalize[k],
						dstLang,
						dstLocalize[k],
					)
				} else {
					pterm.Warning.Printfln("translate fail\nkey:%s\nfrom %s\nto %s\n%s",
						parent+"."+k,
						c.SrcLang,
						dstLang,
						srcLocalize[k],
					)
				}
			}

		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64:

			// Numeric types don't need translation, just copy if key doesn't exist
			if _, exists := dstLocalize[k]; !exists {
				log.Infof("copy numeric value [%s]%d from %s to %s", parent+"."+k, x, c.SrcLang, dstLang)

				dstLocalize[k] = v

				if afterTranslate != nil {
					afterTranslate(dstLang, parent+"."+k, candy.ToString(v), candy.ToString(v))
				}

				if tx != nil {
					err := tx.Update(dstLang.Lang, parent+"."+k, srcLocalize[k])
					if err != nil {
						log.Errorf("err:%v", err)
						return err
					}
				}

				pterm.Success.Printfln("key:%s copied numeric value: %v", parent+"."+k, x)
			}

		case map[string]any:
			dstSub, err := toDstSub(k)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			err = translate(parent+"."+k, x, dstLang, dstSub, c, tx)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

		case map[any]any, map[string]string, map[int64]any, map[int64]string,
			map[int]string, map[float64]any, map[float64]string:
			dstSub, err := toDstSub(k)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			srcSub, err := convertToStringAnyMap(x)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			err = translate(parent+"."+k, srcSub, dstLang, dstSub, c, tx)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

		default:
			log.Errorf("unsupported value type: %T for key %s", x, parent+"."+k)
			return errors.New("unsupported value type: " + fmt.Sprintf("%T", x))
		}
	}

	return nil
}

// LoadLocalize loads localization data from a file
func LoadLocalize(path string, localizer Localizer) (map[string]any, error) {
	m := make(map[string]any)
	err := localizer.UnmarshalFromFile(path, &m)
	if err != nil {
		if os.IsNotExist(err) {
			return m, nil
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	return m, nil
}

// SaveLocalize saves localization data to a file
func SaveLocalize(path string, m map[string]any, localizer Localizer) error {
	err := localizer.MarshalToFile(path, m)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
