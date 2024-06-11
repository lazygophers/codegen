package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var i18nCmd = &cobra.Command{
	Use: "i18n",
	RunE: func(cmd *cobra.Command, args []string) error {
		state.Config.State.I18n = true
		return runGenI18n(cmd, args)
	},
}

func runGenI18n(c *cobra.Command, args []string) (err error) {
	err = codegen.GenerateI18n(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func initI18n() {
	i18nCmd.Short = state.Localize(state.I18nTagCliGenStateI18nShort)
	i18nCmd.Long = state.Localize(state.I18nTagCliGenStateI18nLong)

	stateCmd.AddCommand(i18nCmd)
}
