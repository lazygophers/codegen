package i18n

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/anyx"
	"github.com/pterm/pterm"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type TransacteConfig struct {
	SrcFile            string
	SrcLang            *Language
	Langs              []*Language
	OverwriteKeyPrefix []string
	Overwrite          bool

	Localizer  Localizer
	Translator Translator
}

func Translate(c *TransacteConfig) error {
	srcLocalize, err := LoadLocalize(c.SrcFile, c.Localizer)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
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

		err = translate("", srcLocalize, dstLang, dstLocalize, c)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = saveLocalizer(filepath.Join(filepath.Dir(c.SrcFile), dstLang.Lang+filepath.Ext(c.SrcFile)), dstLocalize, c.Localizer)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return nil
}

func translate(parent string, srcLocalize map[string]any, dstLang *Language, dstLocalize map[string]any, c *TransacteConfig) (err error) {
	toDstSub := func(k string) map[string]any {
		sub := make(map[string]any)
		if v, ok := dstLocalize[k]; ok {
			switch x := v.(type) {
			case map[string]any:
				sub = x

			case map[string]string:
				for k, v := range x {
					sub[k] = v
				}

			case map[any]any:
				for k, v := range x {
					sub[anyx.ToString(k)] = v
				}

			case map[int64]any:
				for k, v := range x {
					sub[strconv.FormatInt(k, 10)] = v
				}

			case map[int64]string:
				for k, v := range x {
					sub[strconv.FormatInt(k, 10)] = v
				}

			case map[float64]any:
				sub = map[string]any{}
				for k, v := range x {
					sub[strconv.FormatFloat(k, 'f', -1, 64)] = v
				}

			case map[float64]string:
				for k, v := range x {
					sub[strconv.FormatFloat(k, 'f', -1, 64)] = v
				}

			default:
				log.Errorf("data:%v", x)
				log.Panicf("unknown type %T", x)
			}
		}

		dstLocalize[k] = sub
		return sub
	}

	for k, v := range srcLocalize {
		log.Infof("try translate %s from %s to %s", parent+"."+k, c.SrcLang, dstLang)

		switch x := v.(type) {
		case string:
			tran := func() {
				log.Infof("try translate [%s]%s from %s to %s", parent+"."+k, x, c.SrcLang, dstLang)

				traget, err := c.Translator.Translate(c.SrcLang.Tag, dstLang.Tag, x)
				if err != nil {
					log.Errorf("err:%v", err)
					pterm.Warning.Printfln("translate fail\nkey:%s\nfrom %s\nto %s\n%s",
						parent+"."+k,
						c.SrcLang,
						dstLang,
						x,
					)
				} else {
					pterm.Success.Printfln("key:%s\nfrom(%s) %s\nto(%s) %s",
						parent+"."+k,

						c.SrcLang,
						x,

						dstLang,
						traget,
					)

					dstLocalize[k] = traget
				}
			}

			// 找得到的时候，判断一下要不要覆盖
			if _, ok := dstLocalize[k]; !ok {
				tran()
			} else if c.Overwrite {
				tran()
			} else {
				for _, prefix := range c.OverwriteKeyPrefix {
					if !strings.HasPrefix(parent+"."+k, prefix) {
						continue
					}
					tran()
					break
				}
			}

		case map[string]any:
			dstSub := toDstSub(k)

			err = translate(parent+"."+k, x, dstLang, dstSub, c)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

		case map[any]any:
			dstSub := toDstSub(k)
			srcSub := make(map[string]any, len(x))
			for k, v := range x {
				srcSub[anyx.ToString(k)] = v
			}

			err = translate(parent+"."+k, srcSub, dstLang, dstSub, c)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

		case map[string]string:
			dstSub := toDstSub(k)
			srcSub := make(map[string]any, len(x))
			for k, v := range x {
				srcSub[k] = v
			}

			err = translate(parent+"."+k, srcSub, dstLang, dstSub, c)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

		case map[int64]any:
			dstSub := toDstSub(k)
			srcSub := make(map[string]any, len(x))
			for k, v := range x {
				srcSub[strconv.FormatInt(k, 10)] = v
			}

			err = translate(parent+"."+k, srcSub, dstLang, dstSub, c)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

		case map[int64]string:
			dstSub := toDstSub(k)
			srcSub := make(map[string]any, len(x))
			for k, v := range x {
				srcSub[strconv.FormatInt(k, 10)] = v
			}

			err = translate(parent+"."+k, srcSub, dstLang, dstSub, c)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

		case map[int]string:
			dstSub := toDstSub(k)
			srcSub := make(map[string]any, len(x))
			for k, v := range x {
				srcSub[strconv.Itoa(k)] = v
			}

			err = translate(parent+"."+k, srcSub, dstLang, dstSub, c)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

		case map[float64]any:
			dstSub := toDstSub(k)
			srcSub := make(map[string]any, len(x))
			for k, v := range x {
				srcSub[strconv.FormatFloat(k, 'f', -1, 64)] = v
			}

			err = translate(parent+"."+k, srcSub, dstLang, dstSub, c)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

		case map[float64]string:
			dstSub := toDstSub(k)
			srcSub := make(map[string]any, len(x))
			for k, v := range x {
				srcSub[strconv.FormatFloat(k, 'f', -1, 64)] = v
			}

			err = translate(parent+"."+k, srcSub, dstLang, dstSub, c)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

		default:
			log.Panicf("unsupported type: %T", x)
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
