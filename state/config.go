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
	DisableFieldId bool `json:"disable_field_id,omitempty" yaml:"disable_field_id,omitempty" tom:"disable_field_id,omitempty"`
	// 是否禁用自动添加 created_at 字段相关的内容
	DisableFieldCreatedAt bool `json:"disable_field_created_at,omitempty" yaml:"disable_field_created_at,omitempty" toml:"disable_field_created_at,omitempty"`
	// 是否禁用自动添加 updated_at 字段相关的内容
	DisableFieldUpdatedAt bool `json:"disable_field_updated_at,omitempty" yaml:"disable_field_updated_at,omitempty" toml:"disable_field_updated_at,omitempty"`
	// 是否禁用自动添加 deleted_at 字段相关的内容
	DisableFieldDeletedAt bool `json:"disable_field_deleted_at,omitempty" yaml:"disable_field_deleted_at,omitempty" toml:"disable_field_deleted_at,omitempty"`
	// 是否禁用自动添加 gorm tag: column
	DisableGormTagColumn bool `json:"disable_gorm_tag_column,omitempty" yaml:"disable_gorm_tag_column,omitempty" toml:"disable_gorm_tag_column,omitempty"`

	// 对于原始数据，key 为 tag 名。value.key 为字段名，value.value 为 tag 内容
	// 例如: {"gorm": {"id": "column:id;primaryKey;autoIncrement;not null"}}
	// 从初始化后整理为，key 为字段名，value.key 为 tag 名，value.value 为 tag 内容
	// 例如: {"id": {"gorm": "column:id;primaryKey;autoIncrement;not null"}}
	DefaultTag map[string]map[string]string `json:"default_tag,omitempty" yaml:"default_tag,omitempty" toml:"default_tag,omitempty"`
}

func (p *CfgTables) apply() (err error) {
	if p.Disable {
		p.DisableFieldId = true

		p.DisableFieldCreatedAt = true
		p.DisableFieldUpdatedAt = true
		p.DisableFieldDeletedAt = true

		p.DisableGormTagColumn = true
	}

	if len(p.DefaultTag) == 0 {
		p.DefaultTag = make(map[string]map[string]string, 1)
	}

	if _, ok := p.DefaultTag["gorm"]; !ok {
		p.DefaultTag["gorm"] = make(map[string]string)
	}

	if _, ok := p.DefaultTag["gorm"]["id"]; !p.DisableFieldId && !ok {
		p.DefaultTag["gorm"]["id"] = "column:id;primaryKey;autoIncrement;not null"
	}

	if _, ok := p.DefaultTag["gorm"]["created_at"]; !p.DisableFieldCreatedAt && !ok {
		p.DefaultTag["gorm"]["created_at"] = "autoCreateTime;<-:create;column:created_at;not null"
	}

	if _, ok := p.DefaultTag["gorm"]["updated_at"]; !p.DisableFieldUpdatedAt && !ok {
		p.DefaultTag["gorm"]["updated_at"] = "autoUpdateTime;<-:;column:updated_at;not null"
	}

	if _, ok := p.DefaultTag["gorm"]["deleted_at"]; !p.DisableFieldDeletedAt && !ok {
		p.DefaultTag["gorm"]["deleted_at"] = "column:deleted_at;not null"
	}

	{
		newTag := make(map[string]map[string]string)

		for tag, values := range p.DefaultTag {
			for field, value := range values {
				if _, ok := newTag[field]; !ok {
					newTag[field] = make(map[string]string)
				}

				newTag[field][tag] = value
			}
		}

		p.DefaultTag = newTag
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
