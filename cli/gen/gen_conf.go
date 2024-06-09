package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var confCmd = &cobra.Command{
	Use:  "conf",
	RunE: ranGenConf,
}

func ranGenConf(c *cobra.Command, args []string) (err error) {
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
}
