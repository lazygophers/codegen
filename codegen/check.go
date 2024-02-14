package codegen

import (
	"bytes"
	"errors"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/runtime"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
	"strings"
)

func CheckEnvironments() (err error) {
	// NOTE: 检查 protoc 是否存在
	{
		log.Debugf("check protoc:%v", state.Config.ProtocPath)
		cmd := exec.Command(state.Config.ProtocPath, "--version")
		cmd.Dir = runtime.Pwd()
		cmd.Env = os.Environ()

		output, err := cmd.Output()
		if err != nil {
			//goland:noinspection GoTypeAssertionOnErrors
			switch x := err.(type) {
			case *exec.Error:
				if errors.Is(x.Err, exec.ErrNotFound) {
					pterm.Error.Printfln("protoc not found, please download it at `https://github.com/protocolbuffers/protobuf/releases`")
					log.Errorf("err:%v", err)
					return errors.New("protoc not found, please download it at `https://github.com/protocolbuffers/protobuf/releases`")
				}
			default:
				log.Errorf("err:%v", err)
				return err
			}
		}

		output = bytes.TrimSuffix(output, []byte("\n"))

		log.Infof("protoc version:%s", output)
		pterm.Success.Printfln("protoc version:%s", output)
	}

	// NOTE: 检查 protoc-gen-go 是否存在
	{
		// go install github.com/golang/protobuf/protoc-gen-go
		log.Debugf("check protoc-gen-go:%v", state.Config.ProtoGenGoPath)

		cmd := exec.Command(state.Config.ProtoGenGoPath, "--version")
		cmd.Dir = runtime.Pwd()
		cmd.Env = os.Environ()

		_, err = cmd.Output()
		if err != nil {
			//goland:noinspection GoTypeAssertionOnErrors
			switch x := err.(type) {
			case *exec.ExitError:
				switch x.ExitCode() {
				case 1:
					if !strings.Contains(string(x.Stderr), `unknown argument`) {
						log.Errorf("err:%v", err)
						return err
					}
				default:
					log.Errorf("err:%v", err)
					return err
				}
			case *exec.Error:
				if errors.Is(x.Err, exec.ErrNotFound) {
					pterm.Error.Printfln("protoc-gen-go not found, please install it by running `go install github.com/golang/protobuf/protoc-gen-go`")
					log.Errorf("err:%v", err)
					return errors.New("protoc-gen-go not found, please install it by running `go install github.com/golang/protobuf/protoc-gen-go`")
				}

				log.Errorf("err:%v", err)
				return err
			default:
				log.Errorf("err:%v", err)
				return err
			}
		}
	}

	return nil
}

// protoc-gen-go-grpc
// go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
