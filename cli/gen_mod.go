package cli

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var genModCmd = &cobra.Command{
	Use:   "mod",
	Short: "Generate go.mod",
	Long: `Generate go.mod file from proto file.
If go.mod file already exists,skip generate.
Go version will be set with the version in the current environment by go version.
Go module will be set with the go_package in the proto file.
Will run go mod tidy fater generate go.mod file.`,
	RunE: runGenMod,
}

func runGenMod(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateMod(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func init() {
	genCmd.AddCommand(genModCmd)
}
