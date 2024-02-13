package codegen

import (
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GenPbFile(protoFilePath string) error {
	log.Info("gen pb.go with protoc")
	pterm.Info.Printfln("gen pb.go with protoc")

	var args []string

	// NOTE: proto文件所谓文件夹
	args = append(args, "-I", filepath.ToSlash(filepath.Dir(protoFilePath)))

	// NOTE: protoc-gen-go插件
	args = append(args, "--plugin=protoc-gen-go="+state.Config.ProtoGenGoPath)

	// NOTE: 执行的文件夹所在位置
	var execDir string
	if state.Config.OutputPath != "" {
		execDir = filepath.ToSlash(state.Config.OutputPath)
	} else {
		execDir = filepath.ToSlash(filepath.Dir(filepath.Dir(protoFilePath)))
	}
	args = append(args, "--go_out", execDir)

	// NOTE: proto文件
	args = append(args, protoFilePath)

	cmd := exec.Command(state.Config.ProtocPath, args...)
	cmd.Env = os.Environ()
	cmd.Dir = execDir

	log.Infof("exec:%s", strings.Join(cmd.Args, " "))
	log.Infof("exec at:%s", execDir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("err:%v", err)
		pterm.Error.Printfln("exec %s\nerr:%s", pterm.BgCyan.Sprint(strings.Join(cmd.Args, " ")), output)
		return err
	}

	log.Infof("output:%s", output)

	pterm.Success.Printfln("gen protobuf file by protoc success")

	return nil
}
