package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var cmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Generates cmd folder",
	RunE:  runGenCmd,
}

func runGenCmd(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateCmd(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func init() {
	genCmd.AddCommand(cmdCmd)
}
