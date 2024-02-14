package codegen

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
	"path/filepath"
)

func GenerateMod(pb *PbPackage) (err error) {
	modFilePath := filepath.Join(pb.ProjectRoot(), "go.mod")
	if osx.IsFile(modFilePath) {
		pterm.Warning.Printfln("go.mod file already exists,skip generate %s", modFilePath)
		return nil
	}

	// NOTE: 执行 go mo init
	{
		pterm.Info.Printfln("run go mod init")

		cmd := exec.Command("go", "mod", "init", pb.GoPackage())
		cmd.Dir = pb.ProjectRoot()
		cmd.Env = os.Environ()

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		err = cmd.Run()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		pterm.Success.Printfln("generate go.mod file %s", modFilePath)
	}

	{
		pterm.Info.Printfln("run go mod tidy")

		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = pb.ProjectRoot()
		cmd.Env = os.Environ()

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		err = cmd.Run()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		pterm.Success.Printfln("run go mod tidy")
	}

	return nil
}
