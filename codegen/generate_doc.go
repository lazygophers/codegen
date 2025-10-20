package codegen

import (
	"os"
	"path/filepath"

	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/i18n"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
)

func GenerateDoc(pb *PbPackage) (err error) {
	pterm.Info.Println("try generate documentation")

	// 获取文档输出目录
	docsDir := filepath.Join(pb.ProjectRoot(), "docs")
	if !osx.IsDir(docsDir) {
		err = os.MkdirAll(docsDir, 0755)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	modelMessages := candy.Filter(pb.Messages(), func(message *PbMessage) bool {
		return message.IsTable()
	})

	if len(modelMessages) == 0 {
		pterm.Warning.Println("no Model-prefixed messages found")
		return nil
	}

	tpl, err := GetTemplate(TemplateTypeDocsDatabase, i18n.DefaultLanguage())
	if err != nil {
		log.Errorf("failed to get template: %v", err)
		return err
	}

	file, err := os.OpenFile(filepath.Join(docsDir, i18n.Localize(state.I18nTagCliDocFilenameDatabaseDesign)+".md"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, FilePermDefault)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	err = tpl.Execute(file, map[string]interface{}{
		"Models":        modelMessages,
		"PackageName":   pb.PackageName(),
		"ProtoFile":     pb.ProtoFileName(),
		"GoPackage":     pb.GoPackage(),
		"GoPackageName": pb.GoPackageName(),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	pterm.Success.Println("documentation generated successfully")
	return nil
}
