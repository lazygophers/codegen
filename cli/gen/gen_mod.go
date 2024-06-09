package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var modCmd = &cobra.Command{
	Use:  "mod",
	RunE: runGenMod,
}

func runGenMod(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateMod(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func initMod() {
	modCmd.Short = state.Localize(state.I18nTagCliGenModShort)
	modCmd.Long = state.Localize(state.I18nTagCliGenModLong)

	genCmd.AddCommand(modCmd)
}
