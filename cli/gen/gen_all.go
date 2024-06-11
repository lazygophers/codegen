package gen

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

type GenHook func(cmd *cobra.Command, args []string) (err error)

var GenAllHooks = []GenHook{
	runGenPb,
	runGenMod,
	runEditorconfig,
	runGoreleaser,
	runMakefile,
	runGitignote,
	runDockerignote,
	runGolangci,
	runGenCmd,
	runGenState,
	runGenConfFile,
	runGenTable,
	runGenImpl,
}

var allCmd = &cobra.Command{
	Use:     "all",
	Aliases: []string{"a", "all-actions"},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		for _, hook := range GenAllHooks {
			err = hook(cmd, args)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
		}

		return nil
	},
}

func initAll() {
	allCmd.Short = state.Localize(state.I18nTagCliGenAllShort)
	allCmd.Long = state.Localize(state.I18nTagCliGenAllLong)

	initStateFlags(allCmd)

	genCmd.AddCommand(allCmd)
}
