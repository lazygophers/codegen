package state

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/app"
	"github.com/lazygophers/utils/json"
	"github.com/lazygophers/utils/osx"
	"github.com/lazygophers/utils/runtime"
	"github.com/pelletier/go-toml/v2"
	"github.com/pterm/pterm"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Cfg struct {
	ProtocPath     string `json:"protoc_path,omitempty" yaml:"protoc_path,omitempty" toml:"protoc_path,omitempty"`
	ProtoGenGoPath string `json:"protogen_go_path,omitempty" yaml:"protogen_go_path,omitempty" toml:"protogen_go_path,omitempty"`
}

func (p *Cfg) apply() (err error) {
	if p.ProtocPath == "" {
		if runtime.IsWindows() {
			p.ProtocPath = "protoc.exe"
		} else {
			p.ProtocPath = "protoc"
		}
	}

	if p.ProtoGenGoPath == "" {
		if runtime.IsWindows() {
			p.ProtoGenGoPath = "protoc-gen-go.exe"
		} else {
			p.ProtoGenGoPath = "protoc-gen-go"
		}
	}

	// NOTE: 检查 protoc 是否存在
	{
		log.Debugf("check protoc:%v", p.ProtocPath)
		cmd := exec.Command(p.ProtocPath, "--version")
		cmd.Dir = runtime.Pwd()
		cmd.Env = os.Environ()

		output, err := cmd.Output()
		if err != nil {
			//goland:noinspection GoTypeAssertionOnErrors
			switch x := err.(type) {
			case *exec.Error:
				if errors.Is(x.Err, exec.ErrNotFound) {
					pterm.Error.Printfln("protoc-gen-go not found, please install it by running `go install github.com/golang/protobuf/protoc-gen-go`")
					log.Errorf("err:%v", err)
					return errors.New("protoc-gen-go not found, please install it by running `go install github.com/golang/protobuf/protoc-gen-go`")
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
		log.Debugf("check protoc-gen-go:%v", p.ProtoGenGoPath)

		cmd := exec.Command(p.ProtoGenGoPath, "--version")
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

var Config = new(Cfg)

type Unmarshaler func(reader io.Reader, v interface{}) error

var supportedExt = map[string]Unmarshaler{
	"json": func(reader io.Reader, v interface{}) error {
		return json.NewDecoder(reader).Decode(v)
	},
	"toml": func(reader io.Reader, v interface{}) error {
		return toml.NewDecoder(reader).Decode(v)
	},
	"yaml": func(reader io.Reader, v interface{}) error {
		return yaml.NewDecoder(reader).Decode(v)
	},
	"yml": func(reader io.Reader, v interface{}) error {
		return yaml.NewDecoder(reader).Decode(v)
	},
}

func tryFindConfigPath(baseDir string) string {
	for ext := range supportedExt {
		file := filepath.Join(baseDir, "codegen.cfg."+ext)
		if osx.IsFile(file) {
			return file
		}
	}

	return ""
}

func LoadConfig() error {
	pterm.Debug.Printfln("try load config")

	// 依次确认配置文件的位置

	var path string
	// NOTE: 从环境变量中获取
	log.Warnf("Try to load config from environment variable(%s)", pterm.Gray("LAZYGO_CODEGEN_CONFIG_FILE"))
	path = os.Getenv("LAZYGO_CODEGEN_CONFIG_FILE")

	if path != "" && !osx.IsFile(path) {
		log.Debugf("config file not found:%v", path)
		path = ""
	}

	// NOTE: 从系统目录中获取
	if path == "" {
		log.Warnf("Try to load config from system folder(%s)", pterm.Gray(filepath.Join(runtime.UserConfigDir(), app.Organization)))
		path = tryFindConfigPath(filepath.Join(runtime.UserConfigDir(), app.Organization))
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warnf("Config file not found, use default config")
			log.Debugf("config file not found:%v", path)
		} else {
			log.Warnf("Config file open failed, use default config")
			log.Errorf("err:%v", err)
		}
	} else {
		defer file.Close()
		log.Infof("Config file found, use config from %s", path)
		pterm.Success.Printfln("Config file found, use config from %s", path)

		ext := filepath.Ext(path)
		if unmarshaler, ok := supportedExt[ext[1:]]; ok {
			err = unmarshaler(file, &Config)
			if err != nil {
				log.Errorf("err:%v", err)
			}
		} else {
			log.Errorf("unsupported config file format:%v", ext)
			return fmt.Errorf("unsupported config file format:%v", ext)
		}
	}

	err = Config.apply()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
