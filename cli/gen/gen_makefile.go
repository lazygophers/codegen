package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var makefileCmd = &cobra.Command{
	Use:   "makefile",
	Short: "Generate makefile file",
	RunE:  runMakefile,
}

func runMakefile(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateMakefile(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func init() {
	genCmd.AddCommand(makefileCmd)
}
