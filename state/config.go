package state

import (
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
	"path/filepath"
)

type CfgTables struct {
	Disable bool `json:"disable,omitempty" yaml:"disable,omitempty" toml:"disable,omitempty"`

	// 是否禁用自动添加 id 字段相关的内容
	DisableAutoId bool `json:"disable_auto_id,omitempty" yaml:"disable_auto_id,omitempty" toml:"disable_auto_id,omitempty"`
	// 是否禁用自动添加 created_at 字段相关的内容
	DisableAutoCreatedAt bool `json:"disable_auto_created_at,omitempty" yaml:"disable_auto_created_at,omitempty" toml:"disable_auto_created_at,omitempty"`
	// 是否禁用自动添加 updated_at 字段相关的内容
	DisableAutoUpdatedAt bool `json:"disable_auto_updated_at,omitempty" yaml:"disable_auto_updated_at,omitempty" toml:"disable_auto_updated_at,omitempty"`
	// 是否禁用自动添加 deleted_at 字段相关的内容
	DisableAutoDeletedAt bool `json:"disable_auto_deleted_at,omitempty" yaml:"disable_auto_deleted_at,omitempty" toml:"disable_auto_deleted_at,omitempty"`

	// 默认添加 id 字段的 gorm tag
	DefaultGormTagId string `json:"default_gorm_tag_id,omitempty" yaml:"default_gorm_tag_id,omitempty" toml:"default_gorm_tag_id,omitempty"`
	// 默认添加 created_at 字段的 gorm tag
	DefaultGormTagCreatedAt string `json:"default_gorm_tag_created_at,omitempty" yaml:"default_gorm_tag_created_at,omitempty" toml:"default_gorm_tag_created_at,omitempty"`
	// 默认添加 updated_at 字段的 gorm tag
	DefaultGormTagUpdatedAt string `json:"default_gorm_tag_updated_at,omitempty" yaml:"default_gorm_tag_updated_at,omitempty" toml:"default_gorm_tag_updated_at,omitempty"`
	// 默认添加 deleted_at 字段的 gorm tag
	DefaultGormTagDeletedAt string `json:"default_gorm_tag_deleted_at,omitempty" yaml:"default_gorm_tag_deleted_at,omitempty" toml:"default_gorm_tag_deleted_at,omitempty"`
}

func (p *CfgTables) apply() (err error) {
	if p.Disable {
		p.DisableAutoId = true

		p.DisableAutoCreatedAt = true
		p.DisableAutoUpdatedAt = true
		p.DisableAutoDeletedAt = true
	}

	if p.DefaultGormTagId == "" {
		p.DefaultGormTagId = "column:id;primaryKey;autoIncrement;not null"
	}

	if p.DefaultGormTagCreatedAt == "" {
		p.DefaultGormTagCreatedAt = "autoCreateTime;<-:create;column:created_at;not null"
	}

	if p.DefaultGormTagUpdatedAt == "" {
		p.DefaultGormTagUpdatedAt = "autoUpdateTime;<-:;column:updated_at;not null"
	}

	if p.DefaultGormTagDeletedAt == "" {
		p.DefaultGormTagDeletedAt = "column:deleted_at;not null"
	}

	return nil
}

type Cfg struct {
	ProtocPath     string `json:"protoc_path,omitempty" yaml:"protoc_path,omitempty" toml:"protoc_path,omitempty"`
	ProtoGenGoPath string `json:"protogen_go_path,omitempty" yaml:"protogen_go_path,omitempty" toml:"protogen_go_path,omitempty"`

	GoModulePrefix string `json:"go_module_prefix,omitempty" yaml:"go_module_prefix,omitempty" toml:"go_module_prefix,omitempty"`

	OutputPath string `json:"output_path,omitempty" yaml:"output_path,omitempty" toml:"output_path,omitempty"`

	Tables *CfgTables `json:"tables,omitempty" yaml:"tables,omitempty" toml:"tables,omitempty"`
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

	if p.Tables == nil {
		p.Tables = new(CfgTables)
	}

	err = p.Tables.apply()
	if err != nil {
		return err
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
