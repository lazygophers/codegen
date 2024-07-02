package gen

import (
	"github.com/lazygophers/codegen/cli/utils"
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
	mergeImplFlage(cmd)

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

	err = codegen.GenerateClient(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func mergeImplFlage(cmd *cobra.Command) {
	if utils.Changed("template-impl-route", cmd) {
		state.Config.Template.Impl.Route = utils.GetString("template-impl-route", cmd)
	}

	if utils.Changed("template-impl-path", cmd) {
		state.Config.Template.Impl.Path = utils.GetString("template-impl-path", cmd)
	}
}

func initImpl() {
	implCmd.Short = state.Localize(state.I18nTagCliGenImplShort)
	implCmd.Long = state.Localize(state.I18nTagCliGenImplLong)

	implCmd.Flags().String("template-impl-route", "", state.Localize(state.I18nTagCliGenImplFlagsTemplateImplRoute))

	implCmd.Flags().String("template-impl-path", "", state.Localize(state.I18nTagCliGenImplFlagsTemplateImplPath))

	genCmd.AddCommand(implCmd)
}
