package codegen

import (
	"embed"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/anyx"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/stringx"
	"github.com/pterm/pterm"
	"os"
	"strings"
	"text/template"
)

//go:embed template/*
//go:embed template/state/*
var embedFs embed.FS

type TemplateType uint8

const (
	TemplateTypeEditorconfig TemplateType = iota + 1
	TemplateTypeOrm

	TemplateTypeStateTable
	TemplateTypeStateConf
	TemplateTypeStateCache
	TemplateTypeStateState
)

func GetTemplate(t TemplateType) (tpl *template.Template, err error) {
	var systemPath, embedPath string

	switch t {
	case TemplateTypeEditorconfig:
		systemPath = state.Config.Template.Editorconfig
		embedPath = "template/.editorconfig"

	case TemplateTypeOrm:
		systemPath = state.Config.Template.Orm
		embedPath = "template/orm.gtpl"

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

var DefaultTemplateFunc = template.FuncMap{
	"ToCamel": stringx.ToCamel,
	"ToSnake": stringx.ToSnake,
	"ToLower": strings.ToLower,
	"ToUpper": strings.ToUpper,
	"ToTitle": strings.ToTitle,

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
}
