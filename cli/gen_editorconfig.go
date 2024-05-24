package cli

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var genEditorconfigCmd = &cobra.Command{
	Use:   "editorconfig",
	Short: "Generate .editorconfig file",
	RunE:  runEditorconfig,
}

func runEditorconfig(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateEditorconfig(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func init() {
	genCmd.AddCommand(genEditorconfigCmd)
}
