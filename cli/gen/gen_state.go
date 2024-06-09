package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var GenStateHook = []GenHook{
	ranGenConf,
	ranGenCache,
	runGenTable,
}

var stateCmd = &cobra.Command{
	Use:  "state",
	RunE: runGenState,
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

func initState() {
	stateCmd.Short = state.Localize(state.I18nTagCliGenStateShort)
	stateCmd.Long = state.Localize(state.I18nTagCliGenStateLong)

	genCmd.AddCommand(stateCmd)
}
