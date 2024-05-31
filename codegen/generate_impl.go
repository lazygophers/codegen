package codegen

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"io/fs"
	"os"
	"path/filepath"
)

func initImplDirectory(pb *PbPackage) error {
	if osx.IsDir(GetPath(PathTypeImpl, pb)) {
		return nil
	}

	err := os.MkdirAll(GetPath(PathTypeImpl, pb), fs.ModePerm)
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	return nil
}

func initImplFile(pb *PbPackage, genTo string) (err error) {
	p := filepath.Join(GetPath(PathTypeImpl, pb), genTo+".go")
	if osx.IsFile(p) {
		return nil
	}

	tpl, err := GetTemplate(TemplateTypeImpl)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	file, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(0666))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	err = tpl.Execute(file, map[string]interface{}{
		"PB": pb,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	_, err = file.WriteString("\n")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func generateImpl(pb *PbPackage, rpc *PbRPC) (err error) {
	pterm.Info.Printfln("try generate impl %s to %s", rpc.Name, rpc.genOption.GenTo)

	err = initImplFile(pb, rpc.genOption.GenTo)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	path := filepath.Join(GetPath(PathTypeImpl, pb), rpc.genOption.GenTo+".go")

	log.Infof("gen impl action %s", rpc.genOption.Action)

	// TODO: 缓存
	goFile, err := ParseGoFile(path)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	for _, goFunc := range goFile.funcs {
		if goFunc.RecvType != "" {
			continue
		}

		if goFunc.Name == rpc.Name {
			pterm.Warning.Printfln("%s is exist, skip generate", rpc.Name)
			return nil
		}
	}

	args := map[string]any{
		"PB":           pb,
		"RpcName":      rpc.Name,
		"RequestType":  rpc.rpc.RequestType,
		"ResponseType": rpc.rpc.ReturnsType,
	}

	if rpc.genOption.Model != "" {
		log.Infof("model:%s", rpc.genOption.Model)

		args["Model"] = rpc.genOption.Model

		msg := pb.GetMessage(rpc.genOption.Model)
		if msg != nil {
			pkField := msg.PrimaryField()
			if pkField != nil {
				args["PrimaryKey"] = pkField.Name
				args["PrimaryKeyType"] = pkField.Type()
			}
		}
	}

	tpl, err := GetTemplate(TemplateTypeImplAction, rpc.genOption.Action)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, fs.FileMode(0666))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	//_, err = file.Seek(0, io.SeekEnd)
	//if err != nil {
	//	log.Errorf("err:%v", err)
	//	return err
	//}

	err = tpl.Execute(file, args)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	_, err = file.WriteString("\n")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	pterm.Success.Printfln("generate %s success", rpc.Name)

	return nil
}

func GenerateImpl(pb *PbPackage) (err error) {
	pterm.Info.Printfln("Generating impl directory...")

	err = initImplDirectory(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Info("————————————————————分割线————————————————")

	for _, rpc := range pb.RPCs() {
		err = generateImpl(pb, rpc)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		log.Info(rpc.Name)
		log.Info(rpc.options)
		log.Info(rpc.comment)
		log.Info(rpc.genOption)
	}

	return nil
}
