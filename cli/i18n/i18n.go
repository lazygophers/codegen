package i18n

import (
	"github.com/lazygophers/codegen/state"
	"github.com/spf13/cobra"
)

var i18nCmd = &cobra.Command{
	Use:     "i18n",
	Aliases: []string{"i", "polyglot", "multilanguage", "internationalize"},
}

func Load(rootCmd *cobra.Command) {
	i18nCmd.Short = state.Localize(state.I18nTagCliI18NShort)
	i18nCmd.Long = state.Localize(state.I18nTagCliI18NLong)

	rootCmd.AddCommand(i18nCmd)

	initTran()
	initSupport()
}
