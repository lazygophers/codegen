package codegen

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"io/fs"
	"os"
)

func GenerateConf(pb *PbPackage) (err error) {
	err = initStateDirectory(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// config 文件为覆盖生成单次生成，不会重复读取
	if osx.IsFile(GetPath(PathTypeStateConf, pb)) {
		if !state.Config.Overwrite {
			pterm.Warning.Printfln("config is already existing, skip generation")
			return nil
		}

		pterm.Warning.Printfln("config is already existing, will overwrite")
	}

	// table 文件为覆盖生成
	args := map[string]interface{}{
		"PB": pb,
		"GoImports": []string{
			pb.GoPackage(),
			"github.com/lazygophers/log",
			"github.com/lazygophers/utils/common",
			"github.com/lazygophers/utils/db",
		},
	}

	tpl, err := GetTemplate(TemplateTypeStateConf)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	file, err := os.OpenFile(GetPath(PathTypeStateConf, pb), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(0666))
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