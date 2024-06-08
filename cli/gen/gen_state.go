package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var GenStateHook = []GenHook{
	ranGenConf,
	ranGenCache,
	runGenTable,
}

var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "Generates state and state folder",
	RunE:  runGenState,
}

func runGenState(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateState(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	for _, hook := range GenStateHook {
		err = hook(cmd, args)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return nil
}

func init() {
	genCmd.AddCommand(stateCmd)
}
