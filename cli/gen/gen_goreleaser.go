package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var goreleaserCmd = &cobra.Command{
	Use:  "goreleaser",
	RunE: runGoreleaser,
}

func runGoreleaser(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateGoreleaser(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func initGoreleaser() {
	goreleaserCmd.Short = state.Localize(state.I18nTagCliGenGoreleaserShort)
	goreleaserCmd.Long = state.Localize(state.I18nTagCliGenGoreleaserLong)

	genCmd.AddCommand(goreleaserCmd)
}
