package codegen

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"io/fs"
	"os"
)

func GenerateGitignote(pb *PbPackage) (err error) {
	editorFile := GetPath(PathTypeGitignore, pb)

	if !state.Config.Overwrite {
		if osx.IsFile(editorFile) {
			pterm.Warning.Printfln(".gitignore file already exists,skip generate %s", editorFile)
			return nil
		}
	}

	pterm.Info.Printfln("update .gitignore")

	tpl, err := GetTemplate(TemplateTypeGitignore)
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

	pterm.Success.Printfln("update .gitignore successfully")

	return nil
}

func GenerateDockerignote(pb *PbPackage) (err error) {
	editorFile := GetPath(PathTypeDockerignore, pb)

	if !state.Config.Overwrite {
		if osx.IsFile(editorFile) {
			pterm.Warning.Printfln(".dockerignore file already exists,skip generate %s", editorFile)
			return nil
		}
	}

	pterm.Info.Printfln("update .dockerignore")

	tpl, err := GetTemplate(TemplateTypeDockerignore)
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

	pterm.Success.Printfln("update .dockerignore successfully")

	return nil
}
