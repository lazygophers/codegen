package codegen

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"os"
)

func GenerateGoreleaser(pb *PbPackage) (err error) {
	// cache 文件为覆盖生成单次生成，不会重复读取
	if osx.IsFile(GetPath(PathTypeGoreleaser, pb)) {
		if !state.Config.Overwrite {
			pterm.Warning.Printfln("main is already existing, skip generation")
			return nil
		}

		pterm.Warning.Printfln("main is already existing, will overwrite")
	}

	// table 文件为覆盖生成
	args := map[string]interface{}{
		"PB": pb,
	}

	tpl, err := GetTemplate(TemplateTypeGoreleaser)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	file, err := os.OpenFile(GetPath(PathTypeGoreleaser, pb), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FilePermDefault)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	err = tpl.Execute(file, args)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
