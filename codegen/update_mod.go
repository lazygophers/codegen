package codegen

import (
	"errors"
	"fmt"
	"github.com/lazygophers/codegen/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/osx"
	"github.com/pterm/pterm"
	"golang.org/x/mod/modfile"
	"io"
	"os"
	"os/exec"
	"regexp"
	"resty.dev/v3"
	"strings"
	"time"
)

func lazyGoModIncudeMatcher(repo string) bool {
	if len(state.LazyConfig.GoMod.IncludeRepo) == 0 {
		return true
	}

	for _, include := range state.LazyConfig.GoMod.IncludeRepo {
		if strings.HasPrefix(repo, include) {
			return true
		}

		matched, err := regexp.MatchString(include, repo)
		if err != nil {
			log.Errorf("err:%v", err)
			continue
		}

		if matched {
			return true
		}
	}

	return false
}

func UpdateGoMod(modPath string) (err error) {
	if !osx.IsFile(modPath) {
		pterm.Error.Printfln("not found file %s", modPath)
		return errors.New("not found file " + modPath)
	}

	// 获取实际使用的 goproxy
	goproxy := state.LazyConfig.GoMod.Goproxy
	if goproxy == "" {
		// 从环境变量获取
		goproxy = os.Getenv("GOPROXY")
	}

	// 从 go env 获取
	if goproxy == "" {
		buffer, err := exec.Command("go", "env", "GOPROXY").CombinedOutput()
		if err == nil {
			goproxy = strings.ReplaceAll(strings.TrimSpace(string(buffer)), "\n", "")
		} else {
			log.Errorf("err:%v", err)
		}
	}

	if strings.Contains(goproxy, ",") {
		goproxy = goproxy[:strings.Index(goproxy, ",")]
	}

	if goproxy == "" {
		goproxy = "https://proxy.golang.org"
	}

	pterm.Info.Printfln("use goproxy:%s", goproxy)

	client := resty.New().
		SetTimeout(time.Second * 5).
		SetRetryCount(3).
		SetRetryMaxWaitTime(time.Second * 30).
		SetRetryWaitTime(time.Second).
		AddRetryConditions(func(response *resty.Response, err error) bool {
			if err != nil {
				return true
			}

			switch response.StatusCode() {
			case 200:
				return false
			case 429:
				pterm.Warning.Printfln("query go.mod too quickly, wait 1m and retry")
				time.Sleep(time.Minute)
				return true

			default:
				return false
			}
		}).
		SetLogger(log.Clone().SetOutput(io.Discard).SetLevel(log.PanicLevel)).
		SetBaseURL(goproxy)

	if os.Getenv("HTTPS_PROXY") != "" {
		client.SetProxy(os.Getenv("HTTPS_PROXY"))
	} else if os.Getenv("ALL_PROXY") != "" {
		client.SetProxy(os.Getenv("ALL_PROXY"))
	} else if os.Getenv("HTTP_PROXY") != "" {
		client.SetProxy(os.Getenv("HTTP_PROXY"))
	}

	// 解析 git 信息
	var currentBranch string

	gitRepo, err := ParseGit(modPath)
	if err != nil {
		log.Errorf("err:%v", err)
	} else {
		currentBranch = gitRepo.HeadName()
	}

	// 解析 go.mod
	modFileBuf, err := os.ReadFile(modPath)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	mod, err := modfile.Parse(modPath, modFileBuf, nil)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var changed bool

	for _, require := range mod.Require {
		if require.Indirect {
			continue
		}

		if !lazyGoModIncudeMatcher(require.Mod.Path) {
			pterm.Warning.Printfln("%s is skipped", require.Mod.Path)
			continue
		}

		var version string

		// 找一下应该去搜索什么版本
		if !candy.Contains(state.LazyConfig.GoMod.ExcludeBranches, currentBranch) {
			if v, ok := state.LazyConfig.GoMod.Tag[require.Mod.Path]; ok && v != "" {
				version = v
			}
		}

		if v, ok := state.LazyConfig.GoMod.Branch[currentBranch]; ok && v != "" {
			version = v
		}

		if version == "" {
			version = "latest"
		}

		var queryPath string
		switch version {
		case "latest":
			queryPath = fmt.Sprintf("%s/@latest", require.Mod.Path)

		default:
			queryPath = fmt.Sprintf("%s/@v/%v.info", require.Mod.Path, version)

		}

		queryPath = strings.ToLower(queryPath)

		var resp ModVersionRsp
		rsp, err := client.R().SetResult(&resp).Get(queryPath)
		if err != nil {
			pterm.Error.Printfln("find %s version fail", require.Mod.Path)
			log.Errorf("err:%v", err)
			continue
		}

		if rsp.StatusCode() != 200 {
			log.Info(rsp.Request.URL)
			log.Errorf("status code: %d", rsp.StatusCode())
			pterm.Error.Printfln("find %s version fail", require.Mod.Path)
			continue
		}

		if resp.Version == require.Mod.Version {
			pterm.Success.Printfln("%s is already %s, skied", require.Mod.Path, version)
			continue
		}

		pterm.Success.Printfln("update %s\n%s\n%s", require.Mod.Path, require.Mod.Version, resp.Version)

		err = mod.AddRequire(require.Mod.Path, resp.Version)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		changed = true
	}

	if changed {
		modFileBuf, err = mod.Format()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = os.WriteFile(modPath, modFileBuf, 0644)
		if err != nil {
			log.Errorf("err:%v", err)

			if os.IsPermission(err) {
				pterm.Warning.Printfln("write to %s no permissions, try to write to %s.temp", modPath, modPath)
				e := os.WriteFile(modPath+".temp", modFileBuf, 0644)
				if e != nil {
					log.Errorf("err:%v", e)
					return e
				}
			}

			return err
		}
		pterm.Success.Printfln("update go.mod success")

		pterm.Info.Printfln("try exec `go mod tidy` to update go.sum and cache mod")
		err = exec.Command("go", "mod", "tidy").Run()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

	}

	return nil
}

type ModVersionRsp struct {
	Version string `json:"Version"`
}
