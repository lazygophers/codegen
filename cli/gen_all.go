package cli

import (
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

type GenAllHook func(cmd *cobra.Command, args []string) (err error)

var GenAllHooks = []GenAllHook{
	runGenPb,
	runGenMod,
	runEditorconfig,
	ranGenConf,
	runGenTable,
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
