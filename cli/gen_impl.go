package cli

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var genImplCmd = &cobra.Command{
	Use:   "impl",
	Short: "Generate implementation files",
	RunE:  runGenImpl,
}

var runGenImpl = func(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateImpl(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func init() {
	genCmd.AddCommand(genImplCmd)
}
