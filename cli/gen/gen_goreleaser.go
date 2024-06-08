package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/spf13/cobra"
)

var goreleaserCmd = &cobra.Command{
	Use:   "goreleaser",
	Short: "Generate .goreleaser.yaml file",
	RunE:  runGoreleaser,
}

func runGoreleaser(cmd *cobra.Command, args []string) (err error) {
	err = codegen.GenerateGoreleaser(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func init() {
	genCmd.AddCommand(goreleaserCmd)
}
