package codegen

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/pterm/pterm"
)

const (
	// Error messages and installation guides
	protocNotFoundMsg           = "protoc not found"
	protocInstallGuide          = "please download it at `https://github.com/protocolbuffers/protobuf/releases`"
	protocGenGoNotFoundMsg      = "protoc-gen-go not found"
	protocGenGoInstallGuide     = "please install it by running `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`"
	protocGenGoGrpcNotFoundMsg  = "protoc-gen-go-grpc not found"
	protocGenGoGrpcInstallGuide = "please install it by running `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`"
)

// commandCheckConfig holds configuration for checking a command
type commandCheckConfig struct {
	cmdPath            string
	cmdName            string
	notFoundMsg        string
	installGuide       string
	shouldPrintVersion bool
}

// errorMessage returns the formatted error message
func (cfg commandCheckConfig) errorMessage() string {
	return fmt.Sprintf("%s, %s", cfg.notFoundMsg, cfg.installGuide)
}

// checkCommand checks if a command exists and optionally prints its version
func checkCommand(cfg commandCheckConfig) error {
	log.Debugf("check %s:%v", cfg.cmdName, cfg.cmdPath)

	// First, try to locate the command using LookPath
	resolvedPath, err := exec.LookPath(cfg.cmdPath)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			errMsg := cfg.errorMessage()
			pterm.Error.Println(errMsg)
			log.Errorf("err:%v", err)
			return errors.New(errMsg)
		}
		log.Errorf("err:%v", err)
		return err
	}

	if cfg.shouldPrintVersion {
		return printVersion(cfg, resolvedPath)
	}

	log.Infof("%s found at %s", cfg.cmdName, resolvedPath)
	pterm.Success.Printfln("%s found at %s", cfg.cmdName, resolvedPath)
	return nil
}

// printVersion prints the version of a command
func printVersion(cfg commandCheckConfig, cmdPath string) error {
	cmd := exec.Command(cmdPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			// Some commands don't support --version flag
			if exitErr.ExitCode() == 1 && strings.Contains(string(exitErr.Stderr), `unknown argument`) {
				// Command exists but doesn't support --version, that's okay
				log.Infof("%s found at %s (version unavailable)", cfg.cmdName, cmdPath)
				pterm.Success.Printfln("%s found at %s", cfg.cmdName, cmdPath)
				return nil
			}
		}
		log.Errorf("err:%v", err)
		return err
	}

	version := strings.TrimSpace(string(output))
	log.Infof("%s version:%s", cfg.cmdName, version)
	pterm.Success.Printfln("%s version:%s", cfg.cmdName, version)
	return nil
}

// CheckEnvironments checks if required protobuf tools are installed
func CheckEnvironments() error {
	// Check protoc
	err := checkCommand(commandCheckConfig{
		cmdPath:            state.Config.ProtocPath,
		cmdName:            "protoc",
		notFoundMsg:        protocNotFoundMsg,
		installGuide:       protocInstallGuide,
		shouldPrintVersion: true,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// Check protoc-gen-go
	err = checkCommand(commandCheckConfig{
		cmdPath:            state.Config.ProtoGenGoPath,
		cmdName:            "protoc-gen-go",
		notFoundMsg:        protocGenGoNotFoundMsg,
		installGuide:       protocGenGoInstallGuide,
		shouldPrintVersion: true,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// Check protoc-gen-go-grpc
	err = checkCommand(commandCheckConfig{
		cmdPath:            state.Config.ProtoGenGoGrpcPath,
		cmdName:            "protoc-gen-go-grpc",
		notFoundMsg:        protocGenGoGrpcNotFoundMsg,
		installGuide:       protocGenGoGrpcInstallGuide,
		shouldPrintVersion: true,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
