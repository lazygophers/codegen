package cli

import (
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var genAllCmd = &cobra.Command{
	Use:   "all",
	Short: "gen all action",
	Long:  `Generate all action.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = runGenPb(cmd, args)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = runGenMod(cmd, args)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = runEditorconfigfunc(cmd, args)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return nil
	},
}

func init() {
	genCmd.AddCommand(genAllCmd)
}
