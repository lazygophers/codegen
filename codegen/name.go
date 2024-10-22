package codegen

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/stringx"
	"strings"
)

type NameStyleHandler struct {
	Base     func(source string) (target string)
	NamePath func(source string, model string) (path string)
}

var NameStyle = map[state.CfgStyleName]NameStyleHandler{
	state.CfgStyleNameDefault: {
		Base: func(source string) string {
			return stringx.ToSnake(source)
		},
		NamePath: func(source string, modele string) string {
			return source
		},
	},
	state.CfgStyleNameCamel: {
		Base: func(source string) (target string) {
			return stringx.ToSmallCamel(source)
		},
		NamePath: func(source string, modele string) string {
			return stringx.ToSmallCamel(source)
		},
	},
	state.CfgStyleNameSnake: {
		Base: func(source string) (target string) {
			return stringx.ToSnake(source)
		},
		NamePath: func(source string, modele string) string {
			return stringx.ToSnake(source)
		},
	},
	state.CfgStyleNameKebab: {
		Base: func(source string) (target string) {
			return stringx.ToKebab(source)
		},
		NamePath: func(source string, modele string) string {
			return stringx.ToKebab(source)
		},
	},
	state.CfgStyleNamePascal: {
		Base: func(source string) (target string) {
			return stringx.ToCamel(source)
		},
		NamePath: func(source string, modele string) string {
			return stringx.ToCamel(source)
		},
	},

	state.CfgStyleNameSlash: {
		Base: func(source string) (target string) {
			return stringx.ToSlash(source)
		},
		NamePath: func(source string, modele string) string {
			return stringx.ToSlash(source)
		},
	},
	state.CfgStyleNameSlashSnake: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")
			return stringx.ToSlash(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "/" + modele
		},
	},
	state.CfgStyleNameSlashPascal: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")
			return stringx.ToSlash(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "/" + stringx.ToCamel(modele)
		},
	},
	state.CfgStyleNameSlashCamel: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")
			return stringx.ToSlash(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "/" + stringx.ToSmallCamel(modele)
		},
	},
	state.CfgStyleNameSlashKebab: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")
			return stringx.ToSlash(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "/" + stringx.ToKebab(modele)
		},
	},

	state.CfgStyleNameSlashReverse: {
		NamePath: func(source string, modele string) (path string) {
			return strings.Join(candy.Reverse(strings.Split(stringx.ToSnake(source), "_")), "/")
		},
	},
	state.CfgStyleNameSlashReverseSnake: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

			return modele + "/" + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), "/")
		},
	},
	state.CfgStyleNameSlashReversePascal: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

			return stringx.ToCamel(modele) + "/" + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), "/")
		},
	},
	state.CfgStyleNameSlashReverseCamel: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

			return stringx.ToSmallCamel(modele) + "/" + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), "/")
		},
	},
	state.CfgStyleNameSlashReverseKebab: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

			return stringx.ToKebab(modele) + "/" + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), "/")
		},
	},

	state.CfgStyleNameDot: {
		Base: func(source string) (target string) {
			return stringx.ToDot(source)
		},
		NamePath: func(source string, modele string) string {
			return stringx.ToDot(source)
		},
	},
	state.CfgStyleNameDotSnake: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

			return stringx.ToDot(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "." + modele
		},
	},
	state.CfgStyleNameDotPascal: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

			return stringx.ToDot(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "." + stringx.ToCamel(modele)
		},
	},
	state.CfgStyleNameDotCamel: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

			return stringx.ToDot(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "." + stringx.ToSmallCamel(modele)
		},
	},
	state.CfgStyleNameDotKebab: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

			return stringx.ToDot(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele)) + "." + stringx.ToKebab(modele)
		},
	},

	state.CfgStyleNameDotReverse: {
		NamePath: func(source string, modele string) (path string) {
			return strings.Join(candy.Reverse(strings.Split(stringx.ToSnake(source), "_")), ".")
		},
	},
	state.CfgStyleNameDotReverseSnake: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

			return modele + "." + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), ".")
		},
	},
	state.CfgStyleNameDotReversePascal: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

			return stringx.ToCamel(modele) + "." + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), ".")
		},
	},
	state.CfgStyleNameDotReverseCamel: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

			return stringx.ToSmallCamel(modele) + "." + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), ".")
		},
	},
	state.CfgStyleNameDotReverseKebab: {
		NamePath: func(source string, modele string) (path string) {
			modele = strings.TrimPrefix(stringx.ToSnake(modele), "model_")

			return stringx.ToKebab(modele) + "." + strings.Join(candy.Reverse(strings.Split(strings.TrimSuffix(stringx.ToSnake(source), "_"+modele), "_")), ".")
		},
	},
}

func CoverageStyledBase(style state.CfgStyleName, source string) string {
	if v, ok := NameStyle[style]; ok {
		return v.Base(source)
	}

	return NameStyle[state.CfgStyleNameDefault].Base(source)
}

func CoverageStyledNamePath(style state.CfgStyleName, source string, model string) string {
	if v, ok := NameStyle[style]; ok {
		return v.NamePath(source, model)
	}

	return NameStyle[state.CfgStyleNameDefault].NamePath(source, model)
}
