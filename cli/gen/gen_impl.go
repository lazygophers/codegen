package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var implCmd = &cobra.Command{
	Use:  "impl",
	RunE: runGenImpl,
}

var runGenImpl = func(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateImpl(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = codegen.GenerateImplRpcPath(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = codegen.GenerateImplRpcRoute(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func initImpl() {
	implCmd.Short = state.Localize(state.I18nTagCliGenImplShort)
	implCmd.Long = state.Localize(state.I18nTagCliGenImplLong)

	genCmd.AddCommand(implCmd)
}
