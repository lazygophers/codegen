package codegen

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"io/fs"
	"os"
)

func GenerateConf(pb *PbPackage) (err error) {
	err = initStateDirectory(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
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

	// 读取 Models
	{
		var models []string
		candy.Each(pb.Messages(), func(message *PbMessage) {
			if !message.IsTable() {
				return
			}

			log.Infof("find table %s", message.Name)

			models = append(models, message.Name)
		})

		args["Models"] = models
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
