package gen

import (
	"github.com/lazygophers/codegen/cli/utils"
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/i18n"
	"github.com/spf13/cobra"
)

var docCmd = &cobra.Command{
	Use:  "doc",
	RunE: runGenDoc,
}

var runGenDoc = func(cmd *cobra.Command, args []string) (err error) {
	mergeDocFlags(cmd)

	err = codegen.GenerateDoc(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func mergeDocFlags(cmd *cobra.Command) {
	if utils.Changed("doc-language", cmd) {
		state.Config.Language = utils.GetString("doc-language", cmd)
	}

	if utils.Changed("doc-output-dir", cmd) {
		state.Config.DocOutputDir = utils.GetString("doc-output-dir", cmd)
	}
}

func initDoc() {
	docCmd.Short = state.Localize(state.I18nTagCliGenDocShort)
	docCmd.Long = state.Localize(state.I18nTagCliGenDocLong)

	docCmd.Flags().String("doc-language", i18n.DefaultLanguage(), state.Localize(state.I18nTagCliGenDocFlagsLanguage))
	docCmd.Flags().String("doc-output-dir", "docs", state.Localize(state.I18nTagCliGenDocFlagsOutputDir))

	genCmd.AddCommand(docCmd)
}
