package codegen

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"github.com/pterm/pterm"
	"io/fs"
	"os"
)

func GenerateStateTable(pb *PbPackage) (err error) {
	pterm.Info.Printfln("try generate table.cache")

	if !state.Config.State.Table {
		pterm.Warning.Printfln("state.table is disable generation, skipping generation")
		return nil
	}

	err = initStateDirectory(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// table 文件为覆盖生成
	args := map[string]interface{}{
		"PB": pb,
	}

	// 读取 Models
	{
		var models []string
		candy.Each(pb.Messages(), func(message *PbMessage) {
			if !message.IsTable() {
				return
			}

			log.Infof("find table %s", message.FullName)

			models = append(models, message.FullName)
		})

		args["Models"] = models
	}

	// 生成 table.go
	tpl, err := GetTemplate(TemplateTypeStateTable)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// TODO: 自动添加错误码到 proto 文件中

	file, err := os.OpenFile(GetPath(PathTypeStateTable, pb), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(0666))
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

func GenerateOrm(pb *PbPackage) (err error) {
	// table 文件为覆盖生成
	args := map[string]interface{}{
		"PB": pb,
	}

	// 读取 Models
	{
		var models []string
		candy.Each(pb.Messages(), func(message *PbMessage) {
			// 先全部允许。实际使用的时候要考虑被 model 引用的场景
			//if !message.NeedOrm() {
			//	return
			//}

			log.Infof("find orm object %s", message.FullName)

			models = append(models, message.FullName)
		})

		args["Models"] = models
	}

	// 生成 table.go
	tpl, err := GetTemplate(TemplateTypeOrm)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	file, err := os.OpenFile(GetPath(PathTypeOrm, pb), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(0666))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(getFileHeader(pb))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = tpl.Execute(file, args)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func GenerateTableName(pb *PbPackage) (err error) {
	// table 文件为覆盖生成
	args := map[string]interface{}{
		"PB": pb,
	}

	// 读取 Models
	{
		var models []string
		candy.Each(pb.Messages(), func(message *PbMessage) {
			if !message.IsTable() {
				return
			}

			log.Infof("find table %s", message.FullName)

			models = append(models, message.FullName)
		})

		args["Models"] = models
	}

	// 生成 table.go
	tpl, err := GetTemplate(TemplateTypeTableName)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	file, err := os.OpenFile(GetPath(PathTypeTableName, pb), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(0666))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(getFileHeader(pb))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = tpl.Execute(file, args)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func GenerateTableField(pb *PbPackage) (err error) {
	// table 文件为覆盖生成
	args := map[string]interface{}{
		"PB": pb,
	}

	fieldMap := map[string]PbField{}
	candy.Each(pb.Messages(), func(msg *PbMessage) {
		if !msg.IsTable() {
			return
		}

		candy.Each(msg.NormalFields(), func(field *PbNormalField) {
			if old, ok := fieldMap[field.FieldName()]; ok {
				if old.FieldType() != field.FieldType() || old.IsSlice() != field.IsSlice() {
					log.Warnf("field %s field type mismatch", field.FieldName())
					pterm.Warning.Printfln("field %s field type mismatch", field.FieldName())
				}
			} else {
				fieldMap[field.FieldName()] = field
			}
		})
	})
	args["FieldMap"] = fieldMap

	// 生成 table.go
	tpl, err := GetTemplate(TemplateTypeTableField)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	file, err := os.OpenFile(GetPath(PathTypeTableField, pb), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(0666))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	_, err = file.Write(getFileHeader(pb))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = tpl.Execute(file, args)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
