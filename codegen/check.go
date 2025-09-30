package codegen

import (
	"errors"
	"fmt"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/pterm/pterm"
	"os/exec"
	"strings"
)

const (
	// Error messages and installation guides
	protocNotFoundMsg       = "protoc not found"
	protocInstallGuide      = "please download it at `https://github.com/protocolbuffers/protobuf/releases`"
	protocGenGoNotFoundMsg  = "protoc-gen-go not found"
	protocGenGoInstallGuide = "please install it by running `go install github.com/golang/protobuf/protoc-gen-go`"
)

// commandCheckConfig holds configuration for checking a command
type commandCheckConfig struct {
	cmdPath         string
	cmdName         string
	notFoundMsg     string
	installGuide    string
	shouldPrintVersion bool
}

// checkCommand checks if a command exists and optionally prints its version
func checkCommand(cfg commandCheckConfig) error {
	log.Debugf("check %s:%v", cfg.cmdName, cfg.cmdPath)

	// First, try to locate the command using LookPath
	resolvedPath, err := exec.LookPath(cfg.cmdPath)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			errMsg := fmt.Sprintf("%s, %s", cfg.notFoundMsg, cfg.installGuide)
			pterm.Error.Println(errMsg)
			log.Errorf("err:%v", err)
			return errors.New(errMsg)
		}
		log.Errorf("err:%v", err)
		return err
	}

	// Try to get version information
	if cfg.shouldPrintVersion {
		cmd := exec.Command(resolvedPath, "--version")
		output, err := cmd.Output()
		if err != nil {
			//goland:noinspection GoTypeAssertionOnErrors
			switch x := err.(type) {
			case *exec.ExitError:
				// Some commands don't support --version flag
				if x.ExitCode() == 1 && strings.Contains(string(x.Stderr), `unknown argument`) {
					// Command exists but doesn't support --version, that's okay
					log.Infof("%s found at %s (version unavailable)", cfg.cmdName, resolvedPath)
					pterm.Success.Printfln("%s found at %s", cfg.cmdName, resolvedPath)
					return nil
				}
				log.Errorf("err:%v", err)
				return err
			case *exec.Error:
				if errors.Is(x.Err, exec.ErrNotFound) {
					errMsg := fmt.Sprintf("%s, %s", cfg.notFoundMsg, cfg.installGuide)
					pterm.Error.Println(errMsg)
					log.Errorf("err:%v", err)
					return errors.New(errMsg)
				}
				log.Errorf("err:%v", err)
				return err
			default:
				log.Errorf("err:%v", err)
				return err
			}
		}

		version := strings.TrimSpace(string(output))
		log.Infof("%s version:%s", cfg.cmdName, version)
		pterm.Success.Printfln("%s version:%s", cfg.cmdName, version)
	} else {
		log.Infof("%s found at %s", cfg.cmdName, resolvedPath)
		pterm.Success.Printfln("%s found at %s", cfg.cmdName, resolvedPath)
	}

	return nil
}

// CheckEnvironments checks if required protobuf tools are installed
func CheckEnvironments() error {
	// Check protoc
	if err := checkCommand(commandCheckConfig{
		cmdPath:         state.Config.ProtocPath,
		cmdName:         "protoc",
		notFoundMsg:     protocNotFoundMsg,
		installGuide:    protocInstallGuide,
		shouldPrintVersion: true,
	}); err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// Check protoc-gen-go
	if err := checkCommand(commandCheckConfig{
		cmdPath:         state.Config.ProtoGenGoPath,
		cmdName:         "protoc-gen-go",
		notFoundMsg:     protocGenGoNotFoundMsg,
		installGuide:    protocGenGoInstallGuide,
		shouldPrintVersion: true,
	}); err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

// protoc-gen-go-grpc
// go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
