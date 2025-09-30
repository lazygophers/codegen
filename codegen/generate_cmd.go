package codegen

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"io/fs"
	"os"
)

func initCmdDirectory(pb *PbPackage) error {
	if osx.IsDir(GetPath(PathTypeCmd, pb)) {
		return nil
	}

	err := os.MkdirAll(GetPath(PathTypeCmd, pb), fs.ModePerm)
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	return nil
}

func GenerateCmd(pb *PbPackage) (err error) {
	err = initCmdDirectory(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// cache 文件为覆盖生成单次生成，不会重复读取
	if osx.IsFile(GetPath(PathTypeCmdMain, pb)) {
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

	tpl, err := GetTemplate(TemplateTypeCmd)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	file, err := os.OpenFile(GetPath(PathTypeCmdMain, pb), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FilePermDefault)
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
