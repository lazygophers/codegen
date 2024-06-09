package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var editorconfigCmd = &cobra.Command{
	Use:  "editorconfig",
	RunE: runEditorconfig,
}

func runEditorconfig(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateEditorconfig(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func initEditorconfig() {
	editorconfigCmd.Short = state.Localize(state.I18nTagCliGenEditorconfigShort)
	editorconfigCmd.Long = state.Localize(state.I18nTagCliGenEditorconfigLong)

	genCmd.AddCommand(editorconfigCmd)
}
