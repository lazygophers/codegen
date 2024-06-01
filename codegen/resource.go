package codegen

import (
	"embed"
	"fmt"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/anyx"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/osx"
	"github.com/lazygophers/utils/stringx"
	"github.com/pterm/pterm"
	"go.uber.org/atomic"
	"os"
	"strings"
	"text/template"
)

//go:embed template/*
//go:embed template/state/*
//go:embed template/impl/*
var embedFs embed.FS

type TemplateType uint8

const (
	TemplateTypeEditorconfig TemplateType = iota + 1
	TemplateTypeOrm
	TemplateTypeTableName
	TemplateTypeTableField

	TemplateTypeProtoRpc
	TemplateTypeProtoRpcName
	TemplateTypeProtoRpcReq
	TemplateTypeProtoRpcResp

	TemplateTypeStateTable
	TemplateTypeStateConf
	TemplateTypeStateCache
	TemplateTypeStateState

	TemplateTypeImpl
	TemplateTypeImplAction
	TemplateTypeImplPath
)

func GetTemplate(t TemplateType, args ...string) (tpl *template.Template, err error) {
	var systemPath, embedPath string

	switch t {
	case TemplateTypeEditorconfig:
		systemPath = state.Config.Template.Editorconfig
		embedPath = "template/.editorconfig"

	case TemplateTypeOrm:
		systemPath = state.Config.Template.Orm
		embedPath = "template/orm.gtpl"

	case TemplateTypeTableName:
		systemPath = state.Config.Template.TableName
		embedPath = "template/table_name.gtpl"

	case TemplateTypeTableField:
		systemPath = state.Config.Template.TableField
		embedPath = "template/table_field.gen.gtpl"

	case TemplateTypeProtoRpc:
		systemPath = state.Config.Template.Proto.Rpc
		embedPath = "template/proto/rpc_service.gtpl"

	case TemplateTypeProtoRpcName:
		if len(args) != 1 {
			panic("Must provide")
		}
		if v, ok := state.Config.Template.Proto.Action[args[0]]; ok && v != nil {
			systemPath = v.Name
		}
		embedPath = fmt.Sprintf("template/proto/%s.name.rpc.gtpl", args[0])

	case TemplateTypeProtoRpcReq:
		if len(args) != 1 {
			panic("Must provide")
		}
		if v, ok := state.Config.Template.Proto.Action[args[0]]; ok && v != nil {
			systemPath = v.Req
		}
		embedPath = fmt.Sprintf("template/proto/%s.req.rpc.gtpl", args[0])

	case TemplateTypeProtoRpcResp:
		if len(args) != 1 {
			panic("Must provide")
		}
		if v, ok := state.Config.Template.Proto.Action[args[0]]; ok && v != nil {
			systemPath = v.Resp
		}
		embedPath = fmt.Sprintf("template/proto/%s.resp.rpc.gtpl", args[0])

	case TemplateTypeStateTable:
		systemPath = state.Config.Template.Table
		embedPath = "template/state/table.gtpl"

	case TemplateTypeStateConf:
		systemPath = state.Config.Template.Conf
		embedPath = "template/state/config.gtpl"

	case TemplateTypeStateCache:
		systemPath = state.Config.Template.Cache
		embedPath = "template/state/cache.gtpl"

	case TemplateTypeStateState:
		systemPath = state.Config.Template.State
		embedPath = "template/state/state.gtpl"

	case TemplateTypeImpl:
		systemPath = state.Config.Template.Impl.Impl
		embedPath = "template/impl/impl.gtpl"

	case TemplateTypeImplAction:
		systemPath = state.Config.Template.Impl.Action[args[0]]
		if !osx.IsFile(systemPath) {
			systemPath = state.Config.Template.Impl.Action[""]
		}

		embedPath = fmt.Sprintf("template/impl/%s.rpc.gtpl", args[0])
		if !osx.FsHasFile(embedFs, embedPath) {
			embedPath = "template/impl/.rpc.gtpl"
		}

	case TemplateTypeImplPath:
		systemPath = state.Config.Template.Impl.Path
		embedPath = "template/rpc_path.gtpl"

	default:
		panic("unsupported template type")
	}

	tpl = template.New("").Funcs(DefaultTemplateFunc)

	if systemPath != "" {
		log.Infof("useing %s", systemPath)

		buffer, err := os.ReadFile(systemPath)
		if err != nil {
			if os.IsNotExist(err) {
				pterm.Error.Printfln("%s is not found", systemPath)
				return nil, err
			}

			log.Errorf("err:%v", err)
			return nil, err
		}

		tpl, err = tpl.Parse(string(buffer))
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		return tpl, nil
	}

	log.Infof("useing default %s", embedPath)

	buffer, err := embedFs.ReadFile(embedPath)
	if err != nil {
		if os.IsNotExist(err) {
			pterm.Error.Printfln("%s is not found", systemPath)
			return nil, err
		}

		log.Errorf("err:%v", err)
		return nil, err
	}

	tpl, err = tpl.Parse(string(buffer))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return tpl, nil
}

var (
	counters = map[string]*atomic.Int64{}
)

func IncrWithKey(key string, def int64) int64 {
	if v, ok := counters[key]; ok {
		return v.Inc()
	}

	v := atomic.NewInt64(0)
	counters[key] = v

	return def
}

func DecrWithKey(key string, def int64) int64 {
	if v, ok := counters[key]; ok {
		return v.Dec()
	}

	v := atomic.NewInt64(0)
	counters[key] = v

	return def
}

var DefaultTemplateFunc = template.FuncMap{
	"ToCamel":      stringx.ToCamel,
	"ToSmallCamel": stringx.ToSmallCamel,
	"ToSnake":      stringx.ToSnake,
	"ToLower":      strings.ToLower,
	"ToUpper":      strings.ToUpper,
	"ToTitle":      strings.ToTitle,

	"TrimPrefix": strings.TrimPrefix,
	"TrimSuffix": strings.TrimSuffix,
	"HasPrefix":  strings.HasPrefix,
	"HasSuffix":  strings.HasSuffix,

	"PluckString": anyx.PluckString,
	"PluckInt":    anyx.PluckInt,
	"PluckInt32":  anyx.PluckInt32,
	"PluckUint32": anyx.PluckUint32,
	"PluckInt64":  anyx.PluckInt64,
	"PluckUint64": anyx.PluckUint64,

	"SliceEmpty": func(ss []string) bool {
		return len(ss) == 0
	},

	"Unique":   candy.Unique[string],
	"Sort":     candy.Sort[string],
	"Reverse":  candy.Reverse[string],
	"Top":      candy.Top[string],
	"First":    candy.First[string],
	"Last":     candy.Last[string],
	"Contains": candy.Contains[string],

	"IncrKey": IncrWithKey,
	"DecrKey": DecrWithKey,
}
