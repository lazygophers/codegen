package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var makefileCmd = &cobra.Command{
	Use:  "makefile",
	RunE: runMakefile,
}

func runMakefile(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateMakefile(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func initMakefile() {
	makefileCmd.Short = state.Localize(state.I18nTagCliGenMakefileShort)
	makefileCmd.Long = state.Localize(state.I18nTagCliGenMakefileLong)

	genCmd.AddCommand(makefileCmd)
}
