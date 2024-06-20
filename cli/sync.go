package cli

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/codegen/cli/utils"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/utils/app"
	"github.com/lazygophers/utils/cryptox"
	"github.com/lazygophers/utils/json"
	"github.com/lazygophers/utils/osx"
	"github.com/lazygophers/utils/runtime"
	"github.com/pelletier/go-toml/v2"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// 加载远端的配置文件
var syncCmd = &cobra.Command{
	Use:         "sync",
	Aliases:     []string{"remote", "sync-remote", "update", "update-remote"},
	Example:     "",
	Args:        cobra.MaximumNArgs(1),
	Annotations: nil,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// 尝试解析一下 config 的文件路径
		var configPath string

		// NOTE: 从环境变量中获取
		log.Warnf("Try to load config from environment variable(%s)", pterm.Gray("LAZYGO_CODEGEN_CONFIG_FILE"))
		configPath = os.Getenv("LAZYGO_CODEGEN_CONFIG_FILE")

		if configPath != "" && !osx.IsFile(configPath) {
			log.Debugf("config file not found:%v", configPath)
			configPath = ""
		}

		// NOTE: 从系统目录中获取
		if configPath == "" {
			log.Warnf("Try to load config from system folder(%s)", pterm.Gray(filepath.Join(runtime.UserConfigDir(), app.Organization)))
			configPath = state.TryFindConfigPath(filepath.Join(runtime.UserConfigDir(), app.Organization))
		}

		if configPath == "" {
			configPath = filepath.Join(runtime.UserConfigDir(), app.Organization, "codegen.cfg.yaml")
		}

		var c state.Cfg

		// 尝试读取并解析已有数据
		file, err := os.Open(configPath)
		if err != nil {
			if os.IsNotExist(err) {
				log.Warnf("Config file not found, use default config")
				log.Debugf("config file not found:%v", configPath)
			} else {
				log.Warnf("Config file open failed, use default config")
				log.Errorf("err:%v", err)
			}
		} else {
			log.Infof("Config file found, use config from %s", configPath)
			pterm.Success.Printfln("Config file found, use config from %s", configPath)

			ext := filepath.Ext(configPath)
			if unmarshaler, ok := state.SupportedExt[ext[1:]]; ok {
				err = unmarshaler(file, &c)
				defer file.Close()
				if err != nil {
					log.Errorf("err:%v", err)
				}
			} else {
				defer file.Close()
				log.Errorf("unsupported config file format:%v", ext)
				return fmt.Errorf("unsupported config file format:%v", ext)
			}
		}

		if c.Sync == nil {
			c.Sync = new(state.CfgSync)
		}

		// 看一下有没有远端配置
		remoteUrl := c.Sync.Remote
		if len(args) > 0 {
			remoteUrl = args[0]
		}

		if remoteUrl == "" {
			pterm.Error.Printfln("Please enter the remote configuration address in the first entry.")
			return fmt.Errorf("remote url is empty")
		}

		c.Sync.Remote = remoteUrl

		if cmd.Flag("cache-template-path").Changed {
			c.Sync.CacheTemplatePath = utils.GetString("cache-template-path", cmd)
		}

		// 刚刚有没有校验的东西
		{
			username, password := utils.GetString("username", cmd), utils.GetString("password", cmd)
			if username != "" {
				if password != "" {
					c.Sync.Headers["Authorization"] = "Bearer " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
				} else {
					c.Sync.Headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(username))
				}
			} else if password != "" {
				c.Sync.Headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(password))
			}
		}

		err = syncFromRemote(&c)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		// 存储配置

		// 先备份一下已有的
		if osx.IsFile(configPath) {
			err = os.Rename(configPath, configPath+".bak")
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
		}

		// 如果是从环境变量获取的地址
		if os.Getenv("LAZYGO_CODEGEN_CONFIG_FILE") == configPath {
			file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			defer file.Close()

			switch filepath.Ext(configPath) {
			case ".json":
				encoder := json.NewEncoder(file)
				encoder.SetIndent("", "\t")

				err = encoder.Encode(c)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			case ".yaml", ".yml":
				encoder := yaml.NewEncoder(file)
				encoder.SetIndent(4)

				err = encoder.Encode(c)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			case ".toml":
				err = toml.NewEncoder(file).Encode(c)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

			default:
				log.Panicf("unsupported encoding type: %v", filepath.Ext(configPath))
			}

		} else {
			// 存储到当前目录下的文件（按照文件名 hash)
			fileName := filepath.Join(runtime.UserConfigDir(), app.Organization, "codegen", cryptox.Md5(c.Sync.Remote), "codegen.cfg.yaml")

			if !osx.IsDir(filepath.Dir(fileName)) {
				err = os.MkdirAll(filepath.Dir(fileName), 0777)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}
			}

			file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			defer file.Close()

			encoder := yaml.NewEncoder(file)
			encoder.SetIndent(4)

			err = encoder.Encode(c)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			// 建立软连接
			err = os.RemoveAll(configPath)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			err = os.Symlink(fileName, configPath)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
		}

		pterm.Success.Printfln("Config file saved to %s", configPath)

		return nil
	},
}

func syncFromRemote(c *state.Cfg) error {

	client := resty.New().
		SetTimeout(time.Second * 10).
		SetRetryWaitTime(time.Second).
		SetRetryCount(3).
		SetHeaders(map[string]string{
			"User-Agent": fmt.Sprintf("lazygophers/%s (%s)", app.Version, app.Organization),
		}).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			if err != nil {
				return true
			}

			switch response.StatusCode() {
			case 200:
				return false
			case 429:
				wait := time.Duration(response.Request.Attempt) * time.Second * 30
				pterm.Warning.Printfln("get template %s too quick, wait %s before retrying", response.Request.URL, wait)
				log.Warnf("get template %s too quick, wait %s before retrying", response.Request.URL, wait)
				time.Sleep(wait)
				return true

			default:
				return false
			}
		}).
		SetHeaders(c.Sync.Headers).
		SetLogger(log.Clone().SetOutput(io.Discard))

	{
		u, err := url.Parse(c.Sync.Remote)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		u.Path = filepath.Dir(u.Path)
		client.SetBaseURL(u.String())
	}

	if os.Getenv("HTTPS_PROXY") != "" {
		client.SetProxy(os.Getenv("HTTPS_PROXY"))
	} else if os.Getenv("ALL_PROXY") != "" {
		client.SetProxy(os.Getenv("ALL_PROXY"))
	} else if os.Getenv("HTTP_PROXY") != "" {
		client.SetProxy(os.Getenv("HTTP_PROXY"))
	}

	cacheKey := cryptox.Md5(c.Sync.Remote)
	if c.Sync.CacheTemplatePath == "" {
		c.Sync.CacheTemplatePath = filepath.Join(runtime.UserConfigDir(), app.Organization, "codegen", cacheKey, "template")
	}

	if !osx.IsDir(c.Sync.CacheTemplatePath) {
		pterm.Info.Printfln("try to create template directory at %s", c.Sync.CacheTemplatePath)
		err := os.MkdirAll(c.Sync.CacheTemplatePath, fs.ModePerm)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	var nc state.Cfg
	// 获取新的配置
	{
		resp, err := client.R().
			Get(c.Sync.Remote)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		if resp.StatusCode() != http.StatusOK {
			log.Errorf("unexpected status code: %v", resp.StatusCode())
			return xerror.NewError(int32(resp.StatusCode()))
		}

		ext := filepath.Ext(c.Sync.Remote)
		ext = strings.TrimPrefix(ext, ".")
		if ct := resp.Header().Get("Content-Type"); ct != "" {
			// 解析一下 Content-Type
			ct = strings.TrimSpace(ct)

			if strings.Contains(ct, ";") {
				ct = ct[:strings.Index(ct, ";")]
			}

			switch ct {
			case "application/json", "text/json":
				ext = "json"

			case "application/yaml", "text/yaml", "application/yml", "text/yml":
				ext = "yaml"

			case "application/toml", "text/toml":
				ext = "toml"
			}
		}

		if unmarshaler, ok := state.SupportedExt[ext]; ok {
			err = unmarshaler(bytes.NewBuffer(resp.Body()), &nc)
			if err != nil {
				log.Error(string(resp.Body()))
				log.Errorf("err:%v", err)
				return err
			}
		} else {
			log.Errorf("unsupported config file format:%v", ext)
			return fmt.Errorf("unsupported config file format:%v", ext)
		}
	}

	// 合并一下配置，并且同步一下模版文件
	// 采用相对保守的 merge 方式

	// 同步 sync 的内容
	if nc.Sync != nil {
		for k, v := range c.Sync.Headers {
			nc.Sync.Headers[k] = v
		}

		c.Sync = nc.Sync
	}

	// NOTE: proto 相关的数据不同步
	// NOTE: go-module-prefix 相关的数据不同步
	// NOTE: 输出路径不同步 相关的数据不同步
	// NOTE: 语言信息不同步 相关的数据不同步

	// NOTE: 多语言配置覆盖
	c.I18n = nc.I18n

	// NOTE: 默认标签覆盖
	c.DefaultTag = nc.DefaultTag

	// NOTE: 表配置覆盖
	c.Tables = nc.Tables

	// NOTE: go.mod 配置覆盖
	c.GoMod = nc.GoMod

	// NOTE: state 配置覆盖
	c.State = nc.State

	// NOTE: 处理一下模板配置，模板相关配置采用激进的方式合并
	c.Template = nc.Template
	if c.Template != nil {
		// 深度遍历每一个需要缓存的数据
		var deep func(dst reflect.Value) error
		deep = func(dst reflect.Value) error {
			for dst.Kind() == reflect.Ptr {
				dst = dst.Elem()
			}

			if !dst.IsValid() {
				return nil
			}

			if dst.IsZero() {
				return nil
			}

			switch dst.Kind() {
			case reflect.Struct:
				for i := 0; i < dst.NumField(); i++ {
					err := deep(dst.Field(i))
					if err != nil {
						return err
					}
				}

			case reflect.Slice:
				for i := 0; i < dst.Len(); i++ {
					err := deep(dst.Index(i))
					if err != nil {
						return err
					}
				}

			case reflect.Map:
				for _, key := range dst.MapKeys() {
					err := deep(dst.MapIndex(key))
					if err != nil {
						return err
					}
				}

			case reflect.String:
				pterm.Info.Printfln("try get template from %s", dst.String())
				resp, err := client.R().Get(dst.String())
				if err != nil {
					log.Errorf("get template from %s", resp.Request.URL)
					log.Errorf("err:%v", err)
					return err
				}

				if resp.StatusCode() != http.StatusOK {
					log.Errorf("get template from %s", resp.Request.URL)
					log.Errorf("unexpected status code: %v", resp.StatusCode())
					return xerror.NewError(int32(resp.StatusCode()))
				}

				// md5+原始文件名，如果文件没有变更就可以不用写文件了
				fieName := filepath.Join(c.Sync.CacheTemplatePath, cryptox.Md5(resp.Body())+"-"+filepath.Base(dst.String()))

				if !osx.IsFile(fieName) {
					err = os.WriteFile(fieName, resp.Body(), 0666)
					if err != nil {
						log.Errorf("err:%v", err)
						return err
					}
				}

				dst.SetString(fieName)

			default:
				log.Panicf("unsupported type: %v", dst.Kind())

			}

			return nil
		}

		err := deep(reflect.ValueOf(c.Template))
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return nil
}

func initSync() {
	syncCmd.Short = state.Localize(state.I18nTagCliSyncShort)
	syncCmd.Long = state.Localize(state.I18nTagCliSyncLong)

	syncCmd.Flags().String("cache-template-path", "", state.Localize(state.I18nTagCliSyncFlagsTemplatePath))
	syncCmd.Flags().String("username", "", state.Localize(state.I18nTagCliSyncFlagsUsername))
	syncCmd.Flags().String("password", "", state.Localize(state.I18nTagCliSyncFlagsPassword))

	rootCmd.AddCommand(syncCmd)
}
