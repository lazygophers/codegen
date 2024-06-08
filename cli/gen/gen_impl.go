package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var iImplCmd = &cobra.Command{
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

func init() {
	genCmd.AddCommand(iImplCmd)
}
