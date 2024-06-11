package state

import (
	"fmt"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/app"
	"github.com/lazygophers/utils/defaults"
	"github.com/lazygophers/utils/json"
	"github.com/lazygophers/utils/osx"
	"github.com/lazygophers/utils/runtime"
	"github.com/pelletier/go-toml/v2"
	"github.com/pterm/pterm"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type CfgStyleName string

const (
	// 使用 rpc name
	CfgStyleNameDefault CfgStyleName = ""
	// 蛇形 rpc_name
	CfgStyleNameSnake CfgStyleName = "snake"
	// （大）驼峰 RpcName
	CfgStyleNamePascal CfgStyleName = "pascal"
	// 小驼峰 rpcName
	CfgStyleNameCamel CfgStyleName = "camel"
	// 中划线 rpc-name
	CfgStyleNameKebab CfgStyleName = "kebab"

	// 斜杠  /add/user/record
	CfgStyleNameSlash CfgStyleName = "slash"
	// 斜杠-行为模型风格-蛇形 /add/user_record
	CfgStyleNameSlashSnake CfgStyleName = "slash-snake"
	// 斜杠-行为模型风格-大驼峰 /add/UserRecord
	CfgStyleNameSlashPascal CfgStyleName = "slash-pascal"
	// 斜杠-行为模型风格-小驼峰 /add/userRecord
	CfgStyleNameSlashCamel CfgStyleName = "slash-camel"
	// 斜杠-行为模型风格-中划线 /add/user-record
	CfgStyleNameSlashKebab CfgStyleName = "slash-kebab"

	// 斜杠-反转 /record/user/add/
	CfgStyleNameSlashReverse CfgStyleName = "slash-reverse"
	// 斜杠-反转-行为模型风格-蛇形 /user_record/add
	CfgStyleNameSlashReverseSnake CfgStyleName = "slash-reverse-snake"
	// 斜杠-反转-行为模型风格-大驼峰 /UserRecord/add
	CfgStyleNameSlashReversePascal CfgStyleName = "slash-reverse-pascal"
	// 斜杠-反转-行为模型风格-小驼峰 /userRecord/add
	CfgStyleNameSlashReverseCamel CfgStyleName = "slash-reverse-camel"
	// 斜杠-反转-行为模型风格-中划线 /user-record/add
	CfgStyleNameSlashReverseKebab CfgStyleName = "slash-reverse-kebab"

	// 点  add.user.record
	CfgStyleNameDot CfgStyleName = "dot"
	// 点-行为模型风格-蛇形 add.user_record
	CfgStyleNameDotSnake CfgStyleName = "dot-snake"
	// 点-行为模型风格-大驼峰 add.UserRecord
	CfgStyleNameDotPascal CfgStyleName = "dot-pascal"
	// 点-行为模型风格-小驼峰 add.userRecord
	CfgStyleNameDotCamel CfgStyleName = "dot-camel"
	// 点-行为模型风格-中划线 add.user-record
	CfgStyleNameDotKebab CfgStyleName = "dot-kebab"

	// 点-反转
	CfgStyleNameDotReverse CfgStyleName = "dot-reverse"
	// 点-反转-行为模型风格-蛇形 user_record.add
	CfgStyleNameDotReverseSnake CfgStyleName = "dot-reverse-snake"
	// 点-反转-行为模型风格-大驼峰 UserRecord.add
	CfgStyleNameDotReversePascal CfgStyleName = "dot-reverse-pascal"
	// 点-反转-行为模型风格-小驼峰 userRecord.add
	CfgStyleNameDotReverseCamel CfgStyleName = "dot-reverse-camel"
	// 点-反转-行为模型风格-中划线 user-record.add
	CfgStyleNameDotReverseKebab CfgStyleName = "dot-reverse-kebab"
)

func (p CfgStyleName) String() string {
	return string(p)
}

func (p *CfgStyleName) UnmarshalJSON(data []byte) error {
	*p = CfgStyleName(data)
	return nil
}

func (p *CfgStyleName) UnmarshalTOML(data []byte) error {
	*p = CfgStyleName(data)
	return nil
}

func (p *CfgStyleName) UnmarshalYAML(unmarshal func(any) error) error {
	var tp string
	if err := unmarshal(&tp); err != nil {
		return err
	}
	*p = CfgStyleName(tp)
	return nil
}

func (p *CfgStyleName) MarshalJSON() ([]byte, error) {
	return []byte(*p), nil
}

func (p *CfgStyleName) MarshalTOML() ([]byte, error) {
	return []byte(*p), nil
}

func (p *CfgStyleName) MarshalYAML() (interface{}, error) {
	return *p, nil
}

type CfgStyle struct {
	Path CfgStyleName `json:"path,omitempty" yaml:"path,omitempty" toml:"path,omitempty"`
}

type CfgProtoAction struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	Req  string `json:"req,omitempty" yaml:"req,omitempty" toml:"req,omitempty"`
	Resp string `json:"resp,omitempty" yaml:"resp,omitempty" toml:"resp,omitempty"`
}

type CfgProto struct {
	Action map[string]*CfgProtoAction `json:"action,omitempty" yaml:"action,omitempty" toml:"action,omitempty"`

	Rpc string `json:"rpc,omitempty" yaml:"rpc,omitempty" toml:"rpc,omitempty"`
}

type CfgImpl struct {
	Action map[string]string `json:"action,omitempty" yaml:"action,omitempty" toml:"action,omitempty"`

	Impl  string `json:"impl,omitempty" yaml:"impl,omitempty" toml:"impl,omitempty"`
	Path  string `json:"path,omitempty" yaml:"path,omitempty" toml:"path,omitempty"`
	Route string `json:"route,omitempty" yaml:"route,omitempty" toml:"route,omitempty"`
}

type CfgTemplateState struct {
	Table string `json:"table,omitempty" yaml:"table,omitempty" toml:"table,omitempty"`
	Conf  string `json:"conf,omitempty" yaml:"conf,omitempty" toml:"conf,omitempty"`
	Cache string `json:"cache,omitempty" yaml:"cache,omitempty" toml:"cache,omitempty"`
	State string `json:"state,omitempty" yaml:"state,omitempty" toml:"state,omitempty"`
	I18n  string `json:"i18n,omitempty" yaml:"i18n,omitempty" toml:"i18n,omitempty"`
}

type CfgTemplate struct {
	Editorconfig string `json:"editorconfig,omitempty" yaml:"editorconfig,omitempty" toml:"editorconfig,omitempty"`

	Orm        string `json:"orm,omitempty" yaml:"orm,omitempty" toml:"orm,omitempty"`
	TableName  string `json:"table_name,omitempty" yaml:"table_name,omitempty" toml:"table_name,omitempty"`
	TableField string `json:"table_field,omitempty" yaml:"table_field,omitempty" toml:"table_field,omitempty"`

	Proto *CfgProto `json:"proto,omitempty" yaml:"proto,omitempty" toml:"proto,omitempty"`
	Impl  *CfgImpl  `json:"impl,omitempty" yaml:"impl,omitempty" toml:"impl,omitempty"`

	Main string `json:"main,omitempty" yaml:"main,omitempty" toml:"main,omitempty"`

	Conf string `json:"conf,omitempty" yaml:"conf,omitempty" toml:"conf,omitempty"`

	State *CfgTemplateState `json:"state,omitempty" yaml:"state,omitempty" toml:"state,omitempty"`

	Goreleaser   string `json:"goreleaser,omitempty" yaml:"goreleaser,omitempty" toml:"goreleaser,omitempty"`
	Makefile     string `json:"makefile,omitempty" yaml:"makefile,omitempty" toml:"makefile,omitempty"`
	Golangci     string `json:"golangci,omitempty" yaml:"golangci,omitempty" toml:"golangci,omitempty"`
	Gitignore    string `json:"gitignore,omitempty" yaml:"gitignore,omitempty" toml:"gitignore,omitempty"`
	Dockerignore string `json:"dockerignore,omitempty" yaml:"dockerignore,omitempty" toml:"dockerignore,omitempty"`

	I18nConst string `json:"i18n_const,omitempty" yaml:"i18n_const,omitempty" toml:"i18n_const,omitempty"`
}

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
}

func (p *CfgTables) apply() {
	if p.Disable {
		p.DisableFieldId = true

		p.DisableFieldCreatedAt = true
		p.DisableFieldUpdatedAt = true
		p.DisableFieldDeletedAt = true

		p.DisableGormTagColumn = true
	}
}

type CfgI18n struct {
	GenerateConst bool `json:"generate_const,omitempty" yaml:"generate_const,omitempty" toml:"generate_const,omitempty"`

	Languages    []string `json:"languages,omitempty" yaml:"languages,omitempty" toml:"languages,omitempty"`
	AllLanguages bool     `json:"all_languages,omitempty" yaml:"all_languages,omitempty" toml:"all_languages,omitempty"`

	Translater string `json:"translater,omitempty" yaml:"translater,omitempty" toml:"translater,omitempty"`
}

func (p *CfgI18n) apply() {
	if p.Translater == "" {
		p.Translater = "google-free"
	}

	if p.AllLanguages {
		// do nothing
	} else {
		if len(p.Languages) == 0 {
			p.Languages = []string{
				"en",
				"zh",
				"ja",
				"ko",
				"fr",
				"de",
				"it",
				"ru",
			}
		}
	}
}

type CfgState struct {
	Config bool `json:"config,omitempty" yaml:"config,omitempty" toml:"config,omitempty"`
	Table  bool `json:"table,omitempty" yaml:"table,omitempty" toml:"table,omitempty"`
	Cache  bool `json:"cache,omitempty" yaml:"cache,omitempty" toml:"cache,omitempty"`
	I18n   bool `json:"i18n,omitempty" yaml:"i18n,omitempty" toml:"i18n,omitempty"`
}

type Cfg struct {
	ProtocPath     string `json:"protoc_path,omitempty" yaml:"protoc_path,omitempty" toml:"protoc_path,omitempty"`
	ProtoGenGoPath string `json:"protogen_go_path,omitempty" yaml:"protogen_go_path,omitempty" toml:"protogen_go_path,omitempty"`

	ProtoFiles []string `json:"proto_files,omitempty" yaml:"proto_files,omitempty" toml:"proto_files,omitempty"`

	GoModulePrefix string `json:"go_module_prefix,omitempty" yaml:"go_module_prefix,omitempty" toml:"go_module_prefix,omitempty"`

	OutputPath string `json:"output_path,omitempty" yaml:"output_path,omitempty" toml:"output_path,omitempty"`

	Language string `json:"language,omitempty" yaml:"language,omitempty" toml:"language,omitempty"`

	I18n *CfgI18n `json:"i18n,omitempty" yaml:"i18n,omitempty" toml:"i18n,omitempty"`

	Style *CfgStyle `json:"style,omitempty" yaml:"style,omitempty" toml:"style,omitempty"`

	Template *CfgTemplate `json:"template,omitempty" yaml:"template,omitempty" toml:"template,omitempty"`

	// 对于原始数据，key 为 tag 名。value.key 为字段名，value.value 为 tag 内容
	// 例如: {"gorm": {"id": "column:id;primaryKey;autoIncrement;not null"}}
	// 从初始化后整理为，key 为字段名，value.key 为 tag 名，value.value 为 tag 内容
	// 例如: {"id": {"gorm": "column:id;primaryKey;autoIncrement;not null"}}
	DefaultTag map[string]map[string]string `json:"default_tag,omitempty" yaml:"default_tag,omitempty" toml:"default_tag,omitempty"`

	Tables *CfgTables `json:"tables,omitempty" yaml:"tables,omitempty" toml:"tables,omitempty"`

	State *CfgState `json:"state,omitempty" yaml:"state,omitempty" toml:"state,omitempty"`

	Overwrite bool `json:"-" yaml:"-" toml:"-"`
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

	if p.GoModulePrefix != "" {
		p.GoModulePrefix = strings.TrimSuffix(p.GoModulePrefix, "/")
	}

	if p.Template == nil {
		p.Template = new(CfgTemplate)
	}
	if p.Template.Proto == nil {
		p.Template.Proto = new(CfgProto)
	}
	if len(p.Template.Proto.Action) == 0 {
		p.Template.Proto.Action = make(map[string]*CfgProtoAction)
	}
	if p.Template.Impl == nil {
		p.Template.Impl = new(CfgImpl)
	}
	if len(p.Template.Impl.Action) == 0 {
		p.Template.Impl.Action = make(map[string]string)
	}
	if p.Template.State == nil {
		p.Template.State = new(CfgTemplateState)
	}

	if p.Tables == nil {
		p.Tables = new(CfgTables)
	}

	if p.Style == nil {
		p.Style = new(CfgStyle)
	}
	err = defaults.SetDefaults(p.Style)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if p.I18n == nil {
		p.I18n = new(CfgI18n)
	}
	p.I18n.apply()

	if p.State == nil {
		p.State = &CfgState{
			Config: true,
			Table:  true,
			Cache:  true,
		}
	}

	// NOTE: struct 标签
	{
		if len(p.DefaultTag) == 0 {
			p.DefaultTag = make(map[string]map[string]string, 1)
		}

		if _, ok := p.DefaultTag["gorm"]; !ok {
			p.DefaultTag["gorm"] = make(map[string]string)
		}

		if _, ok := p.DefaultTag["gorm"]["id"]; !p.Tables.DisableFieldId && !ok {
			p.DefaultTag["gorm"]["id"] = "column:id;primaryKey;autoIncrement;not null"
		} else if p.Tables.DisableFieldId && ok {
			delete(p.DefaultTag["gorm"], "id")
		}

		if _, ok := p.DefaultTag["gorm"]["created_at"]; !p.Tables.DisableFieldCreatedAt && !ok {
			p.DefaultTag["gorm"]["created_at"] = "autoCreateTime;<-:create;column:created_at;not null"
		} else if p.Tables.DisableFieldCreatedAt && ok {
			delete(p.DefaultTag["gorm"], "created_at")
		}

		if _, ok := p.DefaultTag["gorm"]["updated_at"]; !p.Tables.DisableFieldUpdatedAt && !ok {
			p.DefaultTag["gorm"]["updated_at"] = "autoUpdateTime;<-:;column:updated_at;not null"
		} else if p.Tables.DisableFieldUpdatedAt && ok {
			delete(p.DefaultTag["gorm"], "updated_at")
		}

		if _, ok := p.DefaultTag["gorm"]["deleted_at"]; !p.Tables.DisableFieldDeletedAt && !ok {
			p.DefaultTag["gorm"]["deleted_at"] = "column:deleted_at;not null"
		} else if p.Tables.DisableFieldDeletedAt && ok {
			delete(p.DefaultTag["gorm"], "deleted_at")
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
