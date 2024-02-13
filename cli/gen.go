package cli

import (
	"errors"
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/osx"
	"github.com/lazygophers/utils/runtime"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"path/filepath"
)

var pb *codegen.PbPackage

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "gen",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		protoFile := filepath.Join(runtime.Pwd(), getString("input", cmd))
		pterm.Info.Println("proto file:", protoFile)

		// NOTE: 生成pb文件
		err = codegen.GenPbFile(protoFile)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		// NOTE: 检查proto文件是否存在
		{
			if !osx.IsFile(protoFile) {
				log.Errorf("file not found:%s", protoFile)
				pterm.Error.Printfln("proto file not found:%s", protoFile)
				return errors.New("proto file not found")
			}
		}

		pb, err = codegen.ParseProto(protoFile)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return nil
	},
}

func init() {
	genCmd.PersistentFlags().StringP("input", "i", "", "input protobuf file")

	rootCmd.AddCommand(genCmd)
}
