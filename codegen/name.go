package codegen

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/stringx"
	"strings"
)

type NameStyleHandler func(source string, model string) (path string)

var NameStyle = map[state.CfgStyleName]NameStyleHandler{
	state.CfgStyleNameDefault: func(source string, modele string) string {
		return source
	},
	state.CfgStyleNameCamel: func(source string, modele string) string {
		return stringx.ToSmallCamel(source)
	},
	state.CfgStyleNameSnake: func(source string, modele string) string {
		return stringx.ToSnake(source)
	},
	state.CfgStyleNameKebab: func(source string, modele string) string {
		return stringx.ToKebab(source)
	},
	state.CfgStyleNamePascal: func(source string, modele string) string {
		return stringx.ToCamel(source)
	},

	state.CfgStyleNameSlash: func(source string, modele string) string {
		return stringx.ToSlash(source)
	},
	state.CfgStyleNameSlashSnake: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")
		return stringx.ToSlash(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "/" + modele
	},
	state.CfgStyleNameSlashPascal: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")
		return stringx.ToSlash(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "/" + stringx.ToCamel(modele)
	},
	state.CfgStyleNameSlashCamel: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")
		return stringx.ToSlash(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "/" + stringx.ToSmallCamel(modele)
	},
	state.CfgStyleNameSlashKebab: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")
		return stringx.ToSlash(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "/" + stringx.ToKebab(modele)
	},

	state.CfgStyleNameSlashReverse: func(source string, modele string) (path string) {
		return strings.Join(candy.Reverse(strings.Split(stringx.ToSnake(source), "_")), "/")
	},
	state.CfgStyleNameSlashReverseSnake: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

		return modele + "/" + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), "/")
	},
	state.CfgStyleNameSlashReversePascal: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

		return stringx.ToCamel(modele) + "/" + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), "/")
	},
	state.CfgStyleNameSlashReverseCamel: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

		return stringx.ToSmallCamel(modele) + "/" + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), "/")
	},
	state.CfgStyleNameSlashReverseKebab: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

		return stringx.ToKebab(modele) + "/" + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), "/")
	},

	state.CfgStyleNameDot: func(source string, modele string) string {
		return stringx.ToDot(source)
	},
	state.CfgStyleNameDotSnake: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

		return stringx.ToDot(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "." + modele
	},
	state.CfgStyleNameDotPascal: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

		return stringx.ToDot(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "." + stringx.ToCamel(modele)
	},
	state.CfgStyleNameDotCamel: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

		return stringx.ToDot(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "." + stringx.ToSmallCamel(modele)
	},
	state.CfgStyleNameDotKebab: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

		return stringx.ToDot(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "." + stringx.ToKebab(modele)
	},

	state.CfgStyleNameDotReverse: func(source string, modele string) (path string) {
		return strings.Join(candy.Reverse(strings.Split(stringx.ToSnake(source), "_")), ".")
	},
	state.CfgStyleNameDotReverseSnake: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

		return modele + "." + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), ".")
	},
	state.CfgStyleNameDotReversePascal: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

		return stringx.ToCamel(modele) + "." + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), ".")
	},
	state.CfgStyleNameDotReverseCamel: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

		return stringx.ToSmallCamel(modele) + "." + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), ".")
	},
	state.CfgStyleNameDotReverseKebab: func(source string, modele string) (path string) {
		modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

		return stringx.ToKebab(modele) + "." + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), ".")
	},
}

func CoverageStyledName(source string, model string) string {
	if v, ok := NameStyle[state.Config.Style.Path]; ok {
		return v(source, model)
	}

	return NameStyle[state.CfgStyleNameDefault](source, model)
}
