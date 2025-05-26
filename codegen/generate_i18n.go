package codegen

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/anyx"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"text/template"
)

func GenerateI18nConst(dstLocalize map[string]any, path string) (err error) {
	pterm.Info.Printfln("try gen i18n const to %s", path)

	localize := map[string]string{}
	var deepKeys func(ps []string, m map[string]interface{})
	deepKeys = func(ps []string, m map[string]interface{}) {
		for k, v := range m {
			key := strings.Join(append(slices.Clone(ps), k), ".")
			localize[key] = ""

			switch x := v.(type) {
			case string:
				localize[key] = x

			case []byte:
				localize[key] = string(x)

			case int:
				localize[key] = strconv.Itoa(x)

			case map[string]any:
				deepKeys(append(slices.Clone(ps), k), x)

			case map[string]string:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[k] = v
				}
				deepKeys(append(slices.Clone(ps), k), sub)

			case map[string][]byte:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[k] = string(v)
				}
				deepKeys(append(slices.Clone(ps), k), sub)

			case map[int64]any:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[strconv.FormatInt(k, 10)] = v
				}
				deepKeys(append(slices.Clone(ps), k), sub)

			case map[int64]string:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[strconv.FormatInt(k, 10)] = v
				}
				deepKeys(append(slices.Clone(ps), k), sub)

			case map[int64][]byte:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[strconv.FormatInt(k, 10)] = string(v)
				}
				deepKeys(append(slices.Clone(ps), k), sub)

			case map[float64]any:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[strconv.FormatFloat(k, 'f', -1, 64)] = v
				}
				deepKeys(append(slices.Clone(ps), k), sub)

			case map[float64]string:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[strconv.FormatFloat(k, 'f', -1, 64)] = v
				}
				deepKeys(append(slices.Clone(ps), k), sub)

			case map[float64][]byte:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[strconv.FormatFloat(k, 'f', -1, 64)] = string(v)
				}
				deepKeys(append(slices.Clone(ps), k), sub)

			case map[any]any:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[anyx.ToString(k)] = v
				}
				deepKeys(append(slices.Clone(ps), k), sub)

			default:
				log.Panicf("unsupported type:%T", x)
			}
		}
	}

	deepKeys([]string{}, dstLocalize)

	tpl, err := GetTemplate(TemplateTypeI18nConst)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(0666))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(getFileHeader(nil))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = tpl.Execute(file, map[string]interface{}{
		"DirName":      filepath.Base(filepath.Dir(path)),
		"Localize":     localize,
		"DeepLocalize": dstLocalize,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

// GenerateI18nField
// 这里是一个比较特殊的东西，我还没有搞明白这么用 go-template 实现，就现用 golang 的方式实现一个
func GenerateI18nField(dstLocalize map[string]any, path string) (err error) {
	pterm.Info.Printfln("try gen i18n field to %s", path)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(0666))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(getFileHeader(nil))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// 写一下 package
	_, err = file.WriteString("package ")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	_, err = file.WriteString(filepath.Base(filepath.Dir(path)))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	_, err = file.WriteString("\n\n")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	_, err = file.WriteString("// Tran Translation function to facilitate the call of generated code\n// NOTE: Need replace\nvar Tran = func(lang string, key string, args ...any) string {\n\tpanic(\"Undefined translation function, please use it after initialization\")\n}\n\n")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	writeLast := func(key string) (err error) {
		err = template.Must(template.New("").Funcs(DefaultTemplateFunc).Parse(I18nFieldLastTemplate)).
			Execute(file, map[string]any{
				"Keys": key,
			})
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		_, err = file.WriteString("\n")
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return nil
	}

	writeMiddle := func(key string, subs []string) (err error) {
		err = template.Must(template.New("").Funcs(DefaultTemplateFunc).Parse(I18NFieldMiddleTemplate)).
			Execute(file, map[string]any{
				"Keys": key,
				"Subs": candy.Sort(subs),
			})
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		_, err = file.WriteString("\n")
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return nil
	}

	writeFirst := func(key string, subs []string) (err error) {
		err = template.Must(template.New("").Funcs(DefaultTemplateFunc).Parse(I18NFieldFirstTemplate)).
			Execute(file, map[string]any{
				"Keys": key,
				"Subs": candy.Sort(subs),
			})
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		_, err = file.WriteString("\n")
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return nil
	}

	// 递归一下
	localize := map[string][]string{}
	var deepKeys func(ps []string, m map[string]interface{}) (err error)
	deepKeys = func(ps []string, m map[string]interface{}) (err error) {
		for k, v := range m {
			key := strings.Join(append(slices.Clone(ps), k), ".")

			switch x := v.(type) {
			case string:
				err = writeLast(key)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			case []byte:
				err = writeLast(key)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			case int:
				err = writeLast(key)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			case map[string]any:
				localize[key] = anyx.MapKeysString(x)

				err = deepKeys(append(slices.Clone(ps), k), x)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			case map[string]string:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[k] = v
				}

				localize[key] = anyx.MapKeysString(sub)

				err = deepKeys(append(slices.Clone(ps), k), sub)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			case map[string][]byte:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[k] = string(v)
				}

				localize[key] = anyx.MapKeysString(sub)

				err = deepKeys(append(slices.Clone(ps), k), sub)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			case map[int64]any:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[strconv.FormatInt(k, 10)] = v
				}

				localize[key] = anyx.MapKeysString(sub)

				err = deepKeys(append(slices.Clone(ps), k), sub)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			case map[int64]string:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[strconv.FormatInt(k, 10)] = v
				}

				localize[key] = anyx.MapKeysString(sub)

				err = deepKeys(append(slices.Clone(ps), k), sub)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			case map[int64][]byte:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[strconv.FormatInt(k, 10)] = string(v)
				}

				localize[key] = anyx.MapKeysString(sub)

				err = deepKeys(append(slices.Clone(ps), k), sub)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			case map[float64]any:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[strconv.FormatFloat(k, 'f', -1, 64)] = v
				}

				localize[key] = anyx.MapKeysString(sub)

				err = deepKeys(append(slices.Clone(ps), k), sub)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			case map[float64]string:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[strconv.FormatFloat(k, 'f', -1, 64)] = v
				}
				err = deepKeys(append(slices.Clone(ps), k), sub)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			case map[float64][]byte:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[strconv.FormatFloat(k, 'f', -1, 64)] = string(v)
				}

				localize[key] = anyx.MapKeysString(sub)

				err = deepKeys(append(slices.Clone(ps), k), sub)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			case map[any]any:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[anyx.ToString(k)] = v
				}

				localize[key] = anyx.MapKeysString(sub)

				err = deepKeys(append(slices.Clone(ps), k), sub)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			default:
				log.Panicf("unsupported type:%T", x)
			}
		}

		return nil
	}

	err = deepKeys([]string{}, dstLocalize)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	for _, value := range candy.Sort(anyx.MapKeysString(localize)) {
		err = writeMiddle(value, localize[value])
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	err = writeFirst("", anyx.MapKeysString(dstLocalize))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func GenerateI18n(pb *PbPackage) (err error) {
	pterm.Info.Printfln("try generate state.i18nL")

	if !state.Config.State.I18n {
		pterm.Warning.Printfln("state.i18nL is disable generation, skipping generation")
		return nil
	}

	err = initStateDirectory(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if osx.IsFile(GetPath(PathTypeStateI18n, pb)) {
		if !state.Config.Overwrite {
			pterm.Warning.Printfln("state.i18n ialready exists, skip generate %s", GetPath(PathTypeStateI18n, pb))
			return nil
		}

		pterm.Warning.Printfln("state.i18n ialready exists, overwrite %s", GetPath(PathTypeStateI18n, pb))
	}

	tpl, err := GetTemplate(TemplateTypeStateI18n)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	file, err := os.OpenFile(GetPath(PathTypeStateI18n, pb), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(0666))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	// TODO: 对于多语言的打包的支持？
	err = tpl.Execute(file, map[string]interface{}{
		"PB": pb,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

var I18nFieldLastTemplate = "type i8nTag{{ ToCamel .Keys }} string\n\nfunc (p i8nTag{{ ToCamel .Keys }}) Key() string {\n    return `{{ .Keys }}`\n}\n\nfunc (p i8nTag{{ ToCamel .Keys }}) String() string {\n    return p.Key()\n}\n\nfunc (p i8nTag{{ ToCamel .Keys }}) Tran(lang string, args ...any) string {\n    return Tran(lang, p.Key(), args...)\n}\n"

var I18NFieldMiddleTemplate = "type i8nTag{{ ToCamel .Keys }} struct { {{ range $key, $value := .Subs}}\n\t// {{ $.Keys }}.{{ $value }}\n\t{{ ToCamel $value }}\ti8nTag{{ ToCamel $.Keys }}{{ ToCamel $value }}{{end}}\n}\n\nfunc (p i8nTag{{ ToCamel .Keys }}) Key() string {\n    return `{{ .Keys }}`\n}\n\nfunc (p i8nTag{{ ToCamel .Keys }}) String() string {\n    return p.Key()\n}\n\nfunc (p i8nTag{{ ToCamel .Keys }}) Tran(lang string, args ...any) string {\n    return Tran(lang, p.Key(), args...)\n}\n\nfunc (p i8nTag{{ ToCamel .Keys }}) TranWithSuffix(lang string, suffix string, args ...any) string {\n    return Tran(lang, p.Key()+\".\"+suffix, args...)\n}\n"

var I18NFieldFirstTemplate = `var I18nTag{{ ToCamel .Keys }} struct { {{ range $key, $value := .Subs}}
	// {{ $value }}
	{{ ToCamel $value }}	i8nTag{{ ToCamel $.Keys }}{{ ToCamel $value }}{{end}}
}
`
