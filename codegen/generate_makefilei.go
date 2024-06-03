package codegen

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"io/fs"
	"os"
)

func GenerateMakefile(pb *PbPackage) (err error) {
	editorFile := GetPath(PathTypeMakefile, pb)

	if !state.Config.Overwrite {
		if osx.IsFile(editorFile) {
			pterm.Warning.Printfln("Makefile file already exists,skip generate %s", editorFile)
			return nil
		}
	}

	pterm.Info.Printfln("update Makefile")

	tpl, err := GetTemplate(TemplateTypeMakefile)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Infof("update %s", editorFile)

	file, err := os.OpenFile(editorFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(0666))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	err = tpl.Execute(file, nil)
	if err != nil {
		if os.IsPermission(err) {
			pterm.Error.Printfln("permission denied")
			return err
		}

		log.Errorf("err:%v", err)
		return err
	}

	pterm.Success.Printfln("update Makefile successfully")

	return nil
}
