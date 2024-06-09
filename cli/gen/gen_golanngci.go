package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var golangciCmd = &cobra.Command{
	Use:     "golang-lint",
	Aliases: []string{"lint"},
	RunE:    runGolangci,
}

func runGolangci(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateGolangci(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func initGolangLint() {
	golangciCmd.Short = state.Localize(state.I18nTagCliGenGolangLintShort)
	golangciCmd.Long = state.Localize(state.I18nTagCliGenGolangLintLong)

	genCmd.AddCommand(golangciCmd)
}
