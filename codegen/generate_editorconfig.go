package codegen

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"os"
	"path/filepath"
)

func GenerateEditorconfig(pb *PbPackage) (err error) {
	editorFile := filepath.Join(pb.ProjectRoot(), ".editorconfig")

	if !state.Config.Overwrite {
		if osx.IsFile(editorFile) {
			pterm.Warning.Printfln(".editorconfig file already exists,skip generate %s", editorFile)
			return nil
		}
	}

	pterm.Info.Printfln("update .editorconfig")

	var buf []byte
	if state.Config.EditorconfigPath != "" {
		log.Infof("useing %s", state.Config.EditorconfigPath)
		buf, err = os.ReadFile(state.Config.EditorconfigPath)
		if err != nil {
			if os.IsNotExist(err) {
				pterm.Error.Printfln("%s is not found", state.Config.EditorconfigPath)
				return err
			}

			log.Errorf("err:%v", err)
			return err
		}
	}

	if len(buf) == 0 {
		log.Info("useing default .editorconfig")
		buf = GetDefaultEditorconfig()
	}

	log.Infof("update %s", editorFile)

	err = os.WriteFile(editorFile, buf, 0600)
	if err != nil {
		if os.IsPermission(err) {
			pterm.Error.Printfln("permission denied")
			return err
		}

		log.Errorf("err:%v", err)
		return err
	}

	pterm.Success.Printfln("update .editorconfig successfully")

	return nil
}
