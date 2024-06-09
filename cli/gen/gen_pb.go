package gen

import (
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var pbCmd = &cobra.Command{
	Use:  "pb",
	RunE: runGenPb,
}

func runGenPb(cmd *cobra.Command, args []string) (err error) {
	pterm.Info.Println("gen protobuf file")

	// NOTE: 生成pb文件
	err = codegen.GenPbFile(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = codegen.GenerateProto(pb)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func initPb() {
	pbCmd.Short = state.Localize(state.I18nTagCliGenPbShort)
	pbCmd.Long = state.Localize(state.I18nTagCliGenPbLong)

	genCmd.AddCommand(pbCmd)
}
