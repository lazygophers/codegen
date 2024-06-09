package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var gitignoreCmd = &cobra.Command{
	Use:  "gitignote",
	RunE: runGitignote,
}

func runGitignote(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateGitignote(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

var dockerignoteCmd = &cobra.Command{
	Use:  "dockerignore",
	RunE: runDockerignote,
}

func runDockerignote(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateDockerignote(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func initIgnote() {
	gitignoreCmd.Short = state.Localize(state.I18nTagCliGenGitignoteShort)
	gitignoreCmd.Long = state.Localize(state.I18nTagCliGenGitignoteLong)

	dockerignoteCmd.Short = state.Localize(state.I18nTagCliGenDockerignoreShort)
	dockerignoteCmd.Long = state.Localize(state.I18nTagCliGenDockerignoreLong)

	genCmd.AddCommand(
		gitignoreCmd,
		dockerignoteCmd,
	)
}
