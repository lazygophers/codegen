package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var golangciCmd = &cobra.Command{
	Use:   "golang-lint",
	Short: "Generate .golangci.yml file",
	RunE:  runGolangci,
}

func runGolangci(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateGolangci(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func init() {
	genCmd.AddCommand(golangciCmd)
}
