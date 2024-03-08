package codegen

import (
	"bytes"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/app"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GenPbFile(pb *PbPackage) error {
	protoFilePath := pb.ProtoFilePath()

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

	// NOTE: 添加注释
	{
		goPbFilePath := filepath.Join(pb.ProjectRoot(), pb.PackageName+".pb.go")

		stat, err := os.Stat(goPbFilePath)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		// NOTE: 读取文件
		goPbFile, err := os.ReadFile(goPbFilePath)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		// NOTE: 添加注释
		idx := strings.Index(string(goPbFile), "\n// source:")
		if idx > 0 {
			var b bytes.Buffer
			b.Write(goPbFile[:idx])
			b.WriteString("\n// \t")
			b.WriteString(app.Name)
			b.WriteString("       v")
			b.WriteString(app.Version)
			b.Write(goPbFile[idx:])

			err = os.WriteFile(goPbFilePath, b.Bytes(), stat.Mode())
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
		}
	}

	pterm.Success.Printfln("gen protobuf file by protoc success")

	return nil
}
