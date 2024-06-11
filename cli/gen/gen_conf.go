package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var statConfCmd = &cobra.Command{
	Use: "conf",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 如果是指定了的情况，不管配置是什么，都强制写需要生成
		state.Config.State.Config = true
		return runGenStateConf(cmd, args)
	},
}

func runGenStateConf(c *cobra.Command, args []string) (err error) {
	err = codegen.GenerateStateConf(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

var confCmd = &cobra.Command{
	Use:  "conf",
	RunE: runGenConfFile,
}

func runGenConfFile(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateConf(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func initConf() {
	confCmd.Short = state.Localize(state.I18nTagCliGenConfShort)
	confCmd.Long = state.Localize(state.I18nTagCliGenConfLong)

	genCmd.AddCommand(confCmd)

	statConfCmd.Short = state.Localize(state.I18nTagCliGenStateConfShort)
	statConfCmd.Long = state.Localize(state.I18nTagCliGenStateConfLong)

	stateCmd.AddCommand(statConfCmd)
}
