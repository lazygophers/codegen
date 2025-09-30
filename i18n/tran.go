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

	err = updateAutoTrain(c.AutoTran, filepath.Dir(c.SrcFile))
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
		}

		err = translate("", srcLocalize, dstLang, dstLocalize, c, tx)
		if err != nil {
			log.Errorf("err:%v", err)
			if tx != nil {
				_ = tx.Rollback()
			}
			return err
		}

		err = saveLocalizer(filepath.Join(filepath.Dir(c.SrcFile), dstLang.Lang+filepath.Ext(c.SrcFile)), dstLocalize, c.Localizer)
		if err != nil {
			log.Errorf("err:%v", err)
			if tx != nil {
				_ = tx.Rollback()
			}

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

// convertToStringAnyMap converts various map types to map[string]any
func convertToStringAnyMap(v any) (map[string]any, error) {
	sub := make(map[string]any)

	switch x := v.(type) {
	case map[string]any:
		return x, nil

	case map[string]string:
		for k, v := range x {
			sub[k] = v
		}

	case map[string]int:
		for k, v := range x {
			sub[k] = v
		}

	case map[any]any:
		for k, v := range x {
			sub[candy.ToString(k)] = v
		}

	case map[any]int:
		for k, v := range x {
			sub[candy.ToString(k)] = v
		}

	case map[int64]any:
		for k, v := range x {
			sub[strconv.FormatInt(k, 10)] = v
		}

	case map[int64]string:
		for k, v := range x {
			sub[strconv.FormatInt(k, 10)] = v
		}

	case map[int64]int:
		for k, v := range x {
			sub[strconv.FormatInt(k, 10)] = v
		}

	case map[float64]any:
		for k, v := range x {
			sub[strconv.FormatFloat(k, 'f', -1, 64)] = v
		}

	case map[float64]string:
		for k, v := range x {
			sub[strconv.FormatFloat(k, 'f', -1, 64)] = v
		}

	case map[int]int:
		for k, v := range x {
			sub[strconv.Itoa(k)] = v
		}

	case map[int]string:
		for k, v := range x {
			sub[strconv.Itoa(k)] = v
		}

	default:
		log.Errorf("unsupported map type %T, data:%v", x, x)
		return nil, errors.New("unsupported map type: " + fmt.Sprintf("%T", x))
	}

	return sub, nil
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

		needTran := func() bool {
			if beforeTranslate != nil {
				if value, ok := dstLocalize[k]; ok {
					if !beforeTranslate(dstLang, parent+"."+k, candy.ToString(v), candy.ToString(value)) {
						return true
					}
				} else {
					if !beforeTranslate(dstLang, parent+"."+k, candy.ToString(v), "") {
						return true
					}
				}
			}

			if c.Overwrite {
				log.Warnf("need tran %s because of overwrite", k)
				return true
			}

			if _, ok := dstLocalize[k]; !ok {
				log.Warnf("need tran %s because of not exist", k)
				return true
			}

			for _, prefix := range c.OverwriteKeyPrefix {
				if !strings.HasPrefix(parent+"."+k, prefix) {
					continue
				}
				log.Warnf("need tran %s because of overwrite key prefix %s", k, prefix)
				return true
			}

			if tx != nil {
				ok, err := tx.NeedTran(dstLang.Lang, parent+"."+k, srcLocalize[k])
				if err != nil {
					log.Errorf("err:%v", err)
				} else if ok {
					log.Warnf("need tran %s because of source text changed", k)
					return true
				}
			}

			return false
		}

		switch x := v.(type) {
		case string:
			tran := func() error {
				log.Infof("try translate [%s]%s from %s to %s", parent+"."+k, x, c.SrcLang, dstLang)

				var targetList []string
				var hasError bool

				for _, s := range strings.Split(x, "\n") {
					if s == "" {
						targetList = append(targetList, s)
						continue
					}

					target, err := c.Translator.Translate(c.SrcLang.Tag, dstLang.Tag, s)
					if err != nil {
						log.Errorf("err:%v", err)
						hasError = true
						break
					} else {

						targetList = append(targetList, target)
					}
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
				return nil
			}

			if needTran() {
				err = tran()
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}
			}

		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64:

			tran := func() error {
				log.Infof("try translate [%s]%d from %s to %s", parent+"."+k, x, c.SrcLang, dstLang)

				pterm.Success.Printfln("key:%s\nfrom(%s) %d\nto(%s) %d",
					parent+"."+k,

					c.SrcLang,
					x,

					dstLang,
					x,
				)

				dstLocalize[k] = candy.ToString(v)

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

				return nil
			}

			if needTran() {
				err = tran()
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}
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

func saveLocalizer(path string, m map[string]any, localizer Localizer) error {
	err := localizer.MarshalToFile(path, m)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
