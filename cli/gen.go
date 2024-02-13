package cli

import (
	"errors"
	"github.com/lazygophers/codegen/codegen"
	"github.com/lazygophers/codegen/state"
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
		// NOTE: 合并输入参数与配置文件
		mergeGenCmdFlags(cmd)

		err = codegen.CheckEnvironments()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		protoFile := filepath.Join(runtime.Pwd(), getString("input", cmd))
		pterm.Info.Println("proto file:", protoFile)

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

func mergeGenCmdFlags(cmd *cobra.Command) {
	if cmd.Flag("protoc").Changed {
		state.Config.ProtocPath = getString("protoc", cmd)
	}

	if cmd.Flag("protoc-gen-go").Changed {
		state.Config.ProtoGenGoPath = getString("protoc-gen-go", cmd)
	}

	if cmd.Flag("output-path").Changed {
		state.Config.OutputPath = getString("output-path", cmd)
	}
}

func init() {
	genCmd.PersistentFlags().String("protoc", "", "protoc path")
	genCmd.PersistentFlags().String("protoc-gen-go", "", "protoc-gen-go path")
	genCmd.PersistentFlags().String("output-path", "", "output path")

	genCmd.PersistentFlags().StringP("input", "i", "", "input protobuf file")

	rootCmd.AddCommand(genCmd)
}
