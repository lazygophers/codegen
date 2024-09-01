package codegen

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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

type ListOptionFieldOption struct {
	Key       string
	KeyField  string
	Value     int
	FieldName string
	FieldType string
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

	args := map[string]any{
		"PB":              pb,
		"RpcName":         rpc.Name,
		"RequestPackage":  pb.GoPackageName(),
		"RequestType":     rpc.rpc.RequestType,
		"ResponsePackage": pb.GoPackageName(),
		"ResponseType":    rpc.rpc.ReturnsType,
	}

	if strings.Contains(rpc.rpc.RequestType, ".") {
		text := rpc.rpc.RequestType
		if idx := strings.LastIndex(text, "."); idx > 0 {
			args["RequestType"] = text[idx+1:]
			text = text[:idx]
		}
		if idx := strings.LastIndex(text, "."); idx > 0 {
			args["RequestPackage"] = text[idx+1:]
		} else {
			args["RequestPackage"] = text
		}
	}

	if strings.Contains(rpc.rpc.ReturnsType, ".") {
		text := rpc.rpc.ReturnsType
		if idx := strings.LastIndex(text, "."); idx > 0 {
			args["ResponseType"] = text[idx+1:]
			text = text[:idx]
		}
		if idx := strings.LastIndex(text, "."); idx > 0 {
			args["ResponsePackage"] = text[idx+1:]
		} else {
			args["ResponsePackage"] = text
		}
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

	// NOTE: 找一下有没有 ListOption
	{
		enum := pb.GetEnum(rpc.rpc.RequestType + "_" + "ListOption")
		if enum != nil {
			var opts []*ListOptionFieldOption

			candy.Each(enum.fields, func(field *PbEnumField) {
				opt := &ListOptionFieldOption{
					Key:      field.field.Name,
					KeyField: GetFullName(field.field),
					Value:    field.field.Integer,
				}

				if opt.Value == 0 {
					return
				}

				if field.comment != nil {
					if v, ok := field.comment.tags["type"]; ok {
						opt.FieldType = strings.ReplaceAll(strings.Join(v.Lines(), ""), " ", "")
					}

					if v, ok := field.comment.tags["field"]; ok {
						opt.FieldName = strings.ReplaceAll(strings.Join(v.Lines(), ""), " ", "")
					}
				}

				// NOTE: 抹平 type 的差异
				switch opt.FieldType {
				case "str", "string":
					opt.FieldType = "string"

				case "strs", "stringslice", "string-slice":
					opt.FieldType = "string-slice"

				case "like":
					opt.FieldType = "like"

				case "llike", "leftlike", "leftLike", "left_like", "left-like":
					opt.FieldType = "left-like"

				case "rlike", "rightlike", "rightLike", "right_like", "right-like":
					opt.FieldType = "right-like"

				case "i", "int":
					opt.FieldType = "int"

				case "i8", "int8":
					opt.FieldType = "int8"

				case "i16", "int16":
					opt.FieldType = "int16"

				case "i32", "int32":
					opt.FieldType = "int32"

				case "i64", "int64":
					opt.FieldType = "int64"

				case "u", "uint":
					opt.FieldType = "uint"

				case "u8", "uint8":
					opt.FieldType = "uint8"

				case "u16", "uint16":
					opt.FieldType = "uint16"

				case "u32", "uint32":
					opt.FieldType = "uint32"

				case "u64", "uint64":
					opt.FieldType = "uint64"

				case "f32", "float32":
					opt.FieldType = "float32"

				case "f64", "float64":
					opt.FieldType = "float64"

				case "b", "bool":
					opt.FieldType = "bool"

				case "is", "ints", "intslice", "int-slice":
					opt.FieldType = "int-slice"

				case "i8s", "int8s", "int8slice", "int8-slice":
					opt.FieldType = "int8-slice"

				case "i16s", "int16s", "int16slice", "int16-slice":
					opt.FieldType = "int16-slice"

				case "i32s", "int32s", "int32slice", "int32-slice":
					opt.FieldType = "int32-slice"

				case "i64s", "int64s", "int64slice", "int64-slice":
					opt.FieldType = "int64-slice"

				case "us", "uints", "uintslice", "uint-slice":
					opt.FieldType = "uint-slice"

				case "u16s", "uint16s", "uint16slice", "uint16-slice":
					opt.FieldType = "uint16-slice"

				case "u32s", "uint32s", "uint32slice", "uint32-slice":
					opt.FieldType = "uint32-slice"

				case "u64s", "uint64s", "uint64slice", "uint64-slice":
					opt.FieldType = "uint64-slice"

				case "f32s", "float32s", "float32slice", "float32-slice":
					opt.FieldType = "float32-slice"

				case "f64s", "float64s", "float64slice", "float64-slice":
					opt.FieldType = "float64-slice"

				case "bs", "bools", "boolslice", "bool-slice":
					opt.FieldType = "bool-slice"

				}

				// NOTE：切割名字
				if opt.FieldName == "" {
					opt.FieldName = strings.TrimPrefix(field.field.Name, "ListOption")
				}

				opts = append(opts, opt)
			})

			args["ListOption"] = opts
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

	goPackage, err := ParseGoDir(GetPath(PathTypeImpl, pb))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if len(goPackage) == 0 {
		goPackage = make(map[string]*GoPackage)
	}

	if goPackage["impl"] == nil {
		goPackage["impl"] = &GoPackage{}
	}

	if len(goPackage["impl"].FuncMap) == 0 {
		goPackage["impl"].FuncMap = make(map[string]*GoFuncNode)
	}

	matchFunc := func(rpc *PbRPC) bool {
		for _, goFuncNode := range goPackage["impl"].FuncMap {
			goFunc := goFuncNode.goFunc

			if goFunc.RecvType != "" {
				continue
			}

			if goFunc.Name == rpc.Name {
				pterm.Warning.Printfln("%s is exist, skip generate", rpc.Name)
				return true
			}
		}

		return false
	}

	for _, rpc := range pb.RPCs() {
		//path := filepath.Join(GetPath(PathTypeImpl, pb), rpc.genOption.GenTo+".go")

		if matchFunc(rpc) {
			continue
		}

		err = generateImpl(pb, rpc)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return nil
}

type RpcRouteOption struct {
	RpcName      string
	Path         string
	Method       string
	Role         string
	SkipGenRoute bool
}

func pbRpc2Route(rpc *PbRPC) *RpcRouteOption {
	opt := &RpcRouteOption{
		RpcName:      rpc.Name,
		Path:         rpc.genOption.Path,
		Method:       rpc.genOption.Method,
		Role:         rpc.genOption.Role,
		SkipGenRoute: rpc.genOption.SkipGenRoute,
	}

	return opt
}

func GenerateImplRpcPath(pb *PbPackage) (err error) {
	pterm.Info.Printfln("Generating impl path")

	// 会覆盖生成
	var rpcs []*RpcRouteOption
	candy.Each(pb.RPCs(), func(rpc *PbRPC) {
		rpcs = append(rpcs, pbRpc2Route(rpc))
	})

	tpl, err := GetTemplate(TemplateTypeImplPath)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	file, err := os.OpenFile(GetPath(PathTypeImplPath, pb), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(0666))
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

	err = tpl.Execute(file, map[string]any{
		"PB":   pb,
		"RPCS": rpcs,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func GenerateImplRpcRoute(pb *PbPackage) (err error) {
	pterm.Info.Printfln("Generating impl routes")

	err = initCmdDirectory(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// 会覆盖生成
	var rpcs []*RpcRouteOption
	candy.Each(pb.RPCs(), func(rpc *PbRPC) {
		rpcs = append(rpcs, pbRpc2Route(rpc))
	})

	tpl, err := GetTemplate(TemplateTypeImplRoute)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	file, err := os.OpenFile(GetPath(PathTypeImplRoute, pb), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(0666))
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

	err = tpl.Execute(file, map[string]any{
		"PB":   pb,
		"RPCS": rpcs,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func generateImplClient(pb *PbPackage, rpc *PbRPC) (err error) {
	pterm.Info.Printfln("try generate impl client %s", rpc.Name)

	args := map[string]any{
		"PB":           pb,
		"RpcName":      rpc.Name,
		"RequestType":  rpc.rpc.RequestType,
		"ResponseType": rpc.rpc.ReturnsType,
		"RPC":          pbRpc2Route(rpc),
	}

	if strings.Contains(rpc.rpc.RequestType, ".") {
		text := rpc.rpc.RequestType
		if idx := strings.LastIndex(text, "."); idx > 0 {
			args["RequestType"] = text[idx+1:]
			text = text[:idx]
		}
		if idx := strings.LastIndex(text, "."); idx > 0 {
			args["RequestPackage"] = text[idx+1:]
		} else {
			args["RequestPackage"] = text
		}
	}

	if strings.Contains(rpc.rpc.ReturnsType, ".") {
		text := rpc.rpc.ReturnsType
		if idx := strings.LastIndex(text, "."); idx > 0 {
			args["ResponseType"] = text[idx+1:]
			text = text[:idx]
		}
		if idx := strings.LastIndex(text, "."); idx > 0 {
			args["ResponsePackage"] = text[idx+1:]
		} else {
			args["ResponsePackage"] = text
		}
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

	tpl, err := GetTemplate(TemplateTypeImplClientCall, rpc.genOption.Action)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	file, err := os.OpenFile(GetPath(PathTypeImplClient, pb), os.O_CREATE|os.O_WRONLY|os.O_APPEND, fs.FileMode(0666))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString("\n")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

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

func initImplClient(pb *PbPackage) error {
	if osx.IsFile(GetPath(PathTypeImplClient, pb)) {
		return nil
	}

	tpl, err := GetTemplate(TemplateTypeImplClient)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	file, err := os.OpenFile(GetPath(PathTypeImplClient, pb), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.FileMode(0666))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer file.Close()

	err = tpl.Execute(file, map[string]any{
		"PB": pb,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func GenerateClient(pb *PbPackage) (err error) {
	pterm.Info.Printfln("Generating impl client...")

	err = initImplDirectory(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	goPackage, err := ParseGoDir(filepath.Dir(GetPath(PathTypeImplClient, pb)))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if len(goPackage) == 0 {
		goPackage = make(map[string]*GoPackage)
	}

	if goPackage[pb.GoPackageName()] == nil {
		goPackage[pb.GoPackageName()] = &GoPackage{}
	}

	if len(goPackage[pb.GoPackageName()].FuncMap) == 0 {
		goPackage[pb.GoPackageName()].FuncMap = make(map[string]*GoFuncNode)
	}

	matchFunc := func(rpc *PbRPC) bool {
		return goPackage[pb.GoPackageName()].FuncMap[rpc.Name] != nil
	}

	for _, rpc := range pb.RPCs() {
		//path := filepath.Join(GetPath(PathTypeImpl, pb), rpc.genOption.GenTo+".go")

		if matchFunc(rpc) {
			pterm.Warning.Printfln("%s is exist, skip generate client", rpc.Name)
			continue
		}

		err = generateImplClient(pb, rpc)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return nil
}
