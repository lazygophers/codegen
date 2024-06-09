package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var cmdCmd = &cobra.Command{
	Use:     "cmd",
	Aliases: []string{"main"},
	RunE:    runGenCmd,
}

func runGenCmd(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateCmd(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func initCmd() {
	cmdCmd.Short = state.Localize(state.I18nTagCliGenCmdShort)
	cmdCmd.Long = state.Localize(state.I18nTagCliGenCmdLong)

	genCmd.AddCommand(cmdCmd)
}
