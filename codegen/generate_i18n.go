package codegen

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

func GenerateI18nConst(dstLocalize map[string]any, path string) (err error) {
	pterm.Info.Printfln("try gen i18n const to %s", path)

	localize := map[string]string{}
	var deepKeys func(ps []string, m map[string]interface{})
	deepKeys = func(ps []string, m map[string]interface{}) {
		for k, v := range m {
			switch x := v.(type) {
			case string:
				localize[strings.Join(append(slices.Clone(ps), k), ".")] = x

			case map[string]any:
				deepKeys(append(slices.Clone(ps), k), x)

			case map[string]string:
				sub := make(map[string]any, len(x))
				for k, v := range x {
					sub[k] = v
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
