package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var gitignoreCmd = &cobra.Command{
	Use:   "gitignote",
	Short: "Generate .gitignote file",
	RunE:  runGitignote,
}

func runGitignote(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateGitignote(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

var genDockerignoteCmd = &cobra.Command{
	Use:   "dockerignore",
	Short: "Generate .dockerignore file",
	RunE:  runDockerignote,
}

func runDockerignote(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateDockerignote(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func init() {
	genCmd.AddCommand(gitignoreCmd, genDockerignoteCmd)
}
