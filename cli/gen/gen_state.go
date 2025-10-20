package gen

import (
	"github.com/lazygophers/codegen/cli/utils"
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var GenStateHook = []GenHook{
	runGenStateConf,
	runGenCache,
	runGenStateTable,
	runGenI18n,
}

var stateCmd = &cobra.Command{
	Use: "state",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mergeStateFlags(cmd)
		return nil
	},
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

func mergeStateFlags(cmd *cobra.Command) {
	if cmd.Flag("table").Changed {
		state.Config.State.Table = utils.GetBool("db", cmd)
	} else if state.LazyConfig.Gen.Table != nil {
		state.Config.State.Table = *state.LazyConfig.Gen.Table
	}

	if cmd.Flag("cache").Changed {
		state.Config.State.Cache = utils.GetBool("cache", cmd)
	} else if state.LazyConfig.Gen.Cache != nil {
		state.Config.State.Cache = *state.LazyConfig.Gen.Cache
	}

	if cmd.Flag("i18n").Changed {
		state.Config.State.I18n = utils.GetBool("i18n", cmd)
	} else if state.LazyConfig.Gen.I18n != nil {
		state.Config.State.I18n = *state.LazyConfig.Gen.I18n
	}

	if cmd.Flag("config").Changed {
		state.Config.State.Config = utils.GetBool("config", cmd)
	} else if state.LazyConfig.Gen.Config != nil {
		state.Config.State.Config = *state.LazyConfig.Gen.Config
	}

	mergeStateTableFlags(cmd)
}

func initStateFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool("table", state.Config.State.Table, state.Localize(state.I18nTagCliGenStateFlagsTable))
	cmd.PersistentFlags().Bool("cache", state.Config.State.Cache, state.Localize(state.I18nTagCliGenStateFlagsCache))
	cmd.PersistentFlags().Bool("i18n", state.Config.State.I18n, state.Localize(state.I18nTagCliGenStateFlagsI18N))
	cmd.PersistentFlags().Bool("config", state.Config.State.Config, state.Localize(state.I18nTagCliGenStateFlagsConfig))
	initStateTableFlags(cmd)
}

func initState() {
	stateCmd.Short = state.Localize(state.I18nTagCliGenStateShort)
	stateCmd.Long = state.Localize(state.I18nTagCliGenStateLong)

	initStateFlags(stateCmd)

	genCmd.AddCommand(stateCmd)
}
