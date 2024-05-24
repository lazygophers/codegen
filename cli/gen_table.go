package cli

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var genTableCmd = &cobra.Command{
	Use:   "table",
	Short: "Generate table",
	Long:  `Generate table from proto file.`,
	RunE:  runGenTable,
}

func runGenTable(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateTable(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func init() {
	genCmd.AddCommand(genTableCmd)
}