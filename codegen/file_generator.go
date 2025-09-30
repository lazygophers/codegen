package codegen

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"os"
)

// FileGeneratorConfig 配置文件生成器的配置
type FileGeneratorConfig struct {
	PathType     PathType
	TemplateType TemplateType
	DisplayName  string
	Args         any
}

// generateConfigFile 生成配置文件的公共函数
func generateConfigFile(pb *PbPackage, cfg FileGeneratorConfig) (err error) {
	filePath := GetPath(cfg.PathType, pb)

	if !state.Config.Overwrite {
		if osx.IsFile(filePath) {
			pterm.Warning.Printfln("%s file already exists,skip generate %s", cfg.DisplayName, filePath)
			return nil
		}
	}

	pterm.Info.Printfln("update %s", cfg.DisplayName)

	tpl, err := GetTemplate(cfg.TemplateType)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Infof("update %s", filePath)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FilePermDefault)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	err = tpl.Execute(file, cfg.Args)
	if err != nil {
		if os.IsPermission(err) {
			pterm.Error.Printfln("permission denied")
			return err
		}

		log.Errorf("err:%v", err)
		return err
	}

	pterm.Success.Printfln("update %s successfully", cfg.DisplayName)

	return nil
}