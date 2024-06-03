package cli

import (
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
	runGolangci,
	runGenCmd,
	runGenState,
	runGenImpl,
}

var genAllCmd = &cobra.Command{
	Use:   "all",
	Short: "gen all action",
	Long:  `Generate all action.`,
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

func init() {
	genCmd.AddCommand(genAllCmd)
}
