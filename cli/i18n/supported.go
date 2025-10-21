package i18n

import (
	"github.com/lazygophers/codegen/cli/utils"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var supportedCmd = &cobra.Command{
	Use:     "support",
	Aliases: []string{"list", "languages", "support-language", "support-languages"},
	RunE: func(cmd *cobra.Command, args []string) error {
		showLanguage := utils.GetString("language", cmd)

		records := [][]string{
			{
				"code",
				"language",
			},
		}

		for _, lang := range allLanguage {
			show := lang.Lang
			if showLanguage != "" {
				show = showLanguage
			}

			records = append(records, []string{
				lang.Lang,
				state.LocalizeWithLanguage(show, state.I18nTagLang+"."+lang.Lang),
			})
		}

		err := pterm.DefaultTable.WithHasHeader(true).
			WithBoxed(true).
			WithData(records).
			Render()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return nil
	},
}

func initSupport() {
	supportedCmd.Short = state.Localize(state.I18nTagCliI18NSupportShort)
	supportedCmd.Long = state.Localize(state.I18nTagCliI18NSupportLong)

	supportedCmd.Flags().String("language", "", state.Localize(state.I18nTagCliI18NSupportFlagsLanguage))

	i18nCmd.AddCommand(supportedCmd)
}
