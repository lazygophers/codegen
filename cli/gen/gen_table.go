package gen

import (
	"github.com/lazygophers/codegen/cli/utils"
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var tableCmd = &cobra.Command{
	Use: "table",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mergeStateTableFlags(cmd)
		return nil
	},
	RunE: runGenTable,
}

func runGenTable(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateOrm(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = codegen.GenerateTableName(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = codegen.GenerateTableField(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return runGenStateTable(cmd, args)
}

var stateTableCmd = &cobra.Command{
	Use: "table",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		mergeStateTableFlags(cmd)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		state.Config.State.Table = true
		return runGenStateTable(cmd, args)
	},
}

func runGenStateTable(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateStateTable(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func mergeStateTableFlags(cmd *cobra.Command) {
	if cmd.Flag("template-state-table").Changed {
		state.Config.Template.State.Table = utils.GetString("template-state-table", cmd)
	}
}

func initStateTableFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("template-state-table", state.Config.Template.State.Table, state.Localize(state.I18nTagCliGenTableFlagsTemplateStateTable))
}

func initTable() {
	tableCmd.Short = state.Localize(state.I18nTagCliGenTableShort)
	tableCmd.Long = state.Localize(state.I18nTagCliGenTableLong)

	stateTableCmd.Short = state.Localize(state.I18nTagCliGenStateTableShort)
	stateTableCmd.Long = state.Localize(state.I18nTagCliGenStateTableLong)

	initStateTableFlags(stateTableCmd)
	initStateTableFlags(tableCmd)

	stateCmd.AddCommand(stateTableCmd)
	genCmd.AddCommand(tableCmd)
}
