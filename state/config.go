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

type CfgSync struct {
	Remote  string            `json:"remote,omitempty" yaml:"remote,omitempty" toml:"remote,omitempty"`
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty" toml:"headers,omitempty"`

	CacheTemplatePath string `json:"cache-template-path,omitempty" yaml:"cache-template-path,omitempty" toml:"cache-template-path,omitempty"`
}

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

	Impl       string `json:"impl,omitempty" yaml:"impl,omitempty" toml:"impl,omitempty"`
	Path       string `json:"path,omitempty" yaml:"path,omitempty" toml:"path,omitempty"`
	Route      string `json:"route,omitempty" yaml:"route,omitempty" toml:"route,omitempty"`
	Client     string `json:"client,omitempty" yaml:"client,omitempty" toml:"client,omitempty"`
	ClientCall string `json:"client_call,omitempty" yaml:"client_call,omitempty" toml:"client_call,omitempty"`
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
	TableName  string `json:"table-name,omitempty" yaml:"table-name,omitempty" toml:"table-name,omitempty"`
	TableField string `json:"table-field,omitempty" yaml:"table-field,omitempty" toml:"table-field,omitempty"`

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

	I18nConst string `json:"i18n-const,omitempty" yaml:"i18n-const,omitempty" toml:"i18n-const,omitempty"`
}

type CfgTables struct {
	Disable bool `json:"disable,omitempty" yaml:"disable,omitempty" toml:"disable,omitempty"`

	// 是否禁用自动添加 id 字段相关的内容
	DisableFieldId bool `json:"disable-field-id,omitempty" yaml:"disable-field-id,omitempty" tom:"disable-field-id,omitempty"`
	// 是否禁用自动添加 created-at 字段相关的内容
	DisableFieldCreatedAt bool `json:"disable-field-created-at,omitempty" yaml:"disable-field-created-at,omitempty" toml:"disable-field-created-at,omitempty"`
	// 是否禁用自动添加 updated-at 字段相关的内容
	DisableFieldUpdatedAt bool `json:"disable-field-updated-at,omitempty" yaml:"disable-field-updated-at,omitempty" toml:"disable-field-updated-at,omitempty"`
	// 是否禁用自动添加 deleted-at 字段相关的内容
	DisableFieldDeletedAt bool `json:"disable-field-deleted-at,omitempty" yaml:"disable-field-deleted-at,omitempty" toml:"disable-field-deleted-at,omitempty"`
	// 是否禁用自动添加 gorm tag: column
	DisableGormTagColumn bool `json:"disable-gorm-tag-column,omitempty" yaml:"disable-gorm-tag-column,omitempty" toml:"disable-gorm-tag-column,omitempty"`
	// 是否禁用自动添加字段类型相关的内容
	DisableFieldType bool `json:"disable-field-type,omitempty" yaml:"disable-field-type,omitempty" tom:"disable-field-type,omitempty"`

	// 是否禁用关于错误：数据未找到的指定错误生成
	DisableErrorNotFound bool `json:"disable_error_not_found,omitempty" toml:"disable_error_not_found,omitempty" yaml:"disable-error-not-found,omitempty"`
	// 是否禁用关于错误：唯一键冲突的指定错误生成
	DisableErrorDuplicateKey bool `json:"disable-error-duplicate-key,omitempty" toml:"disable-error-duplicate-key,omitempty" yaml:"disable-error-duplicate-key,omitempty"`
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
	GenerateConst bool `json:"generate-const,omitempty" yaml:"generate-const,omitempty" toml:"generate-const,omitempty"`
	GenerateField bool `json:"generate-field,omitempty" yaml:"generate-field,omitempty" toml:"generate-field,omitempty"`

	Languages    []string `json:"languages,omitempty" yaml:"languages,omitempty" toml:"languages,omitempty"`
	AllLanguages bool     `json:"all-languages,omitempty" yaml:"all-languages,omitempty" toml:"all-languages,omitempty"`

	Translator string `json:"translator,omitempty" yaml:"translator,omitempty" toml:"translator,omitempty"`
}

func (p *CfgI18n) apply() {
	if p.Translator == "" {
		p.Translator = "google-free"
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

type CfgGoMod struct {
	Goproxy  string `json:"goproxy,omitempty" yaml:"goproxy,omitempty" toml:"goproxy,omitempty"`
	Increase bool   `json:"increase,omitempty" yaml:"increase,omitempty" toml:"increase,omitempty"`
}

type CfgState struct {
	Config bool `json:"config,omitempty" yaml:"config,omitempty" toml:"config,omitempty"`
	Table  bool `json:"table,omitempty" yaml:"table,omitempty" toml:"table,omitempty"`
	Cache  bool `json:"cache,omitempty" yaml:"cache,omitempty" toml:"cache,omitempty"`
	I18n   bool `json:"i18n,omitempty" yaml:"i18n,omitempty" toml:"i18n,omitempty"`
}

type Cfg struct {
	Sync *CfgSync `json:"sync,omitempty" yaml:"sync,omitempty" toml:"sync,omitempty"`

	ProtocPath     string `json:"protoc-path,omitempty" yaml:"protoc-path,omitempty" toml:"protoc-path,omitempty"`
	ProtoGenGoPath string `json:"protogen-go-path,omitempty" yaml:"protogen-go-path,omitempty" toml:"protogen-go-path,omitempty"`

	ProtoFiles []string `json:"proto-files,omitempty" yaml:"proto-files,omitempty" toml:"proto-files,omitempty"`

	GoModulePrefix string `json:"go-module-prefix,omitempty" yaml:"go-module-prefix,omitempty" toml:"go-module-prefix,omitempty"`

	OutputPath string `json:"output-path,omitempty" yaml:"output-path,omitempty" toml:"output-path,omitempty"`

	Language string `json:"language,omitempty" yaml:"language,omitempty" toml:"language,omitempty"`

	I18n *CfgI18n `json:"i18n,omitempty" yaml:"i18n,omitempty" toml:"i18n,omitempty"`

	Style *CfgStyle `json:"style,omitempty" yaml:"style,omitempty" toml:"style,omitempty"`

	Template *CfgTemplate `json:"template,omitempty" yaml:"template,omitempty" toml:"template,omitempty"`

	// 对于原始数据，key 为 tag 名。value.key 为字段名，value.value 为 tag 内容
	// 例如: {"gorm": {"id": "column:id;primaryKey;autoIncrement;not null"}}
	// 从初始化后整理为，key 为字段名，value.key 为 tag 名，value.value 为 tag 内容
	// 例如: {"id": {"gorm": "column:id;primaryKey;autoIncrement;not null"}}
	DefaultTag map[string]map[string]string `json:"default-tag,omitempty" yaml:"default-tag,omitempty" toml:"default-tag,omitempty"`

	Tables *CfgTables `json:"tables,omitempty" yaml:"tables,omitempty" toml:"tables,omitempty"`

	GoMod *CfgGoMod `json:"go-mod,omitempty" yaml:"go-mod,omitempty" toml:"go-mod,omitempty"`

	State *CfgState `json:"state,omitempty" yaml:"state,omitempty" toml:"state,omitempty"`

	Overwrite bool `json:"-" yaml:"-" toml:"-"`
}

func (p *Cfg) apply() (err error) {
	if p.Sync == nil {
		p.Sync = new(CfgSync)
	}

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

		addBoolType := func(typ string) {
			typ = "@" + typ
			if _, ok := p.DefaultTag["gorm"][typ]; !p.Tables.DisableFieldId && !ok {
				p.DefaultTag["gorm"][typ] = "default:0;not null;type:tinyint(1)"
			} else if p.Tables.DisableFieldType && ok {
				delete(p.DefaultTag["gorm"], typ)
			}
		}

		addIntegerType := func(typ string) {
			typ = "@" + typ
			if _, ok := p.DefaultTag["gorm"][typ]; !p.Tables.DisableFieldId && !ok {
				p.DefaultTag["gorm"][typ] = "not null;type:bigint(20)"
			} else if p.Tables.DisableFieldType && ok {
				delete(p.DefaultTag["gorm"], typ)
			}
		}

		addFloatingType := func(typ string) {
			typ = "@" + typ
			if _, ok := p.DefaultTag["gorm"][typ]; !p.Tables.DisableFieldId && !ok {
				p.DefaultTag["gorm"][typ] = "not null;type:decimal(65,8);default:0"
			} else if p.Tables.DisableFieldType && ok {
				delete(p.DefaultTag["gorm"], typ)
			}
		}

		addStringType := func(typ string) {
			typ = "@" + typ
			if _, ok := p.DefaultTag["gorm"][typ]; !p.Tables.DisableFieldId && !ok {
				p.DefaultTag["gorm"][typ] = "type:varchar(255);not null"
			} else if p.Tables.DisableFieldType && ok {
				delete(p.DefaultTag["gorm"], typ)
			}
		}

		addTextType := func(typ string) {
			typ = "@" + typ
			if _, ok := p.DefaultTag["gorm"][typ]; !p.Tables.DisableFieldId && !ok {
				p.DefaultTag["gorm"][typ] = "type:text"
			} else if p.Tables.DisableFieldType && ok {
				delete(p.DefaultTag["gorm"], typ)
			}
		}

		addObjectType := func(typ string) {
			typ = "@" + typ
			if _, ok := p.DefaultTag["gorm"][typ]; !p.Tables.DisableFieldId && !ok {
				p.DefaultTag["gorm"][typ] = "type:json"
			} else if p.Tables.DisableFieldType && ok {
				delete(p.DefaultTag["gorm"], typ)
			}
		}

		addBoolType("bool")

		addIntegerType("int32")
		addIntegerType("int64")
		addIntegerType("uint32")
		addIntegerType("uint64")
		addIntegerType("sint64")
		addIntegerType("sint32")

		addFloatingType("float")
		addFloatingType("double")
		addFloatingType("fixed32")
		addFloatingType("fixed64")

		addStringType("string")
		addStringType("bytes")

		addTextType("array")
		addObjectType("map")
		addObjectType("object")

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

	if p.GoMod == nil {
		p.GoMod = new(CfgGoMod)
	}

	return nil
}

var Config = new(Cfg)

type Unmarshaler func(reader io.Reader, v interface{}) error

var SupportedExt = map[string]Unmarshaler{
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

var supportedExts = [...]string{
	"yaml",
	"yml",
	"json",
	"toml",
}

func TryFindConfigPath(baseDir string) string {
	for _, ext := range supportedExts {
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
		path = TryFindConfigPath(filepath.Join(runtime.UserConfigDir(), app.Organization))
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
		if unmarshaler, ok := SupportedExt[ext[1:]]; ok {
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
