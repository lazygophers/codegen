package cli

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var genPbCmd = &cobra.Command{
	Use:   "pb",
	Short: "gen pb.go",
	Long: `Generate pb.go file from proto file.
Add fields tags to the generated file(pb.go).
`,
	RunE: runGenPb,
}

func runGenPb(cmd *cobra.Command, args []string) (err error) {
	pterm.Info.Println("gen protobuf file")

	err = codegen.GenerateProto(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func init() {
	genCmd.AddCommand(genPbCmd)
}
