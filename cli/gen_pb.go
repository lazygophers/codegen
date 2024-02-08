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
	"os"
	"os/exec"
	"path/filepath"
)

var genPbCmd = &cobra.Command{
	Use:   "pb",
	Short: "pb",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		pterm.Info.Println("gen protobuf file")

		protoFile := filepath.Join(runtime.Pwd(), getString("input", cmd))
		pterm.Info.Println("proto file:", protoFile)

		// NOTE: 生成pb文件
		err = genPbFile(protoFile)
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

		pb, err := codegen.ParseProto(protoFile)
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
	},
}

func init() {
	genPbCmd.Flags().StringP("input", "i", "", "input protobuf file")

	genCmd.AddCommand(genPbCmd)
}

func genPbFile(protoFilePath string) error {
	cmd := exec.Command(state.Config.ProtocPath,
		"-I", filepath.ToSlash(filepath.Dir(protoFilePath)),
		"--plugin=protoc-gen-go="+state.Config.ProtoGenGoPath,
		"--go_out", filepath.ToSlash(filepath.Dir(filepath.Dir(protoFilePath))),
		protoFilePath,
	)
	cmd.Env = os.Environ()
	cmd.Dir = filepath.Dir(filepath.Dir(protoFilePath))

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Infof("output:%s", output)

	pterm.Success.Printfln("gen protobuf file by protoc success")

	return nil
}
