package codegen

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"io/fs"
	"os"
)

func initStateDirectory(pb *PbPackage) error {
	if osx.IsDir(GetPath(PathTypeState, pb)) {
		return nil
	}

	err := os.MkdirAll(GetPath(PathTypeState, pb), fs.ModePerm)
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	return nil
}

func GenerateState(pb *PbPackage) (err error) {
	err = initStateDirectory(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// cache 文件为覆盖生成单次生成，不会重复读取
	if osx.IsFile(GetPath(PathTypeStateState, pb)) {
		if !state.Config.Overwrite {
			pterm.Warning.Printfln("state is already existing, skip generation")
			return nil
		}

		pterm.Warning.Printfln("state is already existing, will overwrite")
	}

	// table 文件为覆盖生成
	args := map[string]interface{}{
		"PB":     pb,
		"Config": state.Config.State,
	}

	tpl, err := GetTemplate(TemplateTypeStateState)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	file, err := os.OpenFile(GetPath(PathTypeStateState, pb), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FilePermDefault)
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
