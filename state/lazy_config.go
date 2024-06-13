package state

import (
	"github.com/lazygophers/log"
	"github.com/pterm/pterm"
	"os"
)

type CfgLazyGoMod struct {
	Goproxy         string            `json:"goproxy,omitempty" yaml:"goproxy,omitempty" toml:"goproxy,omitempty"`
	IncludeRepo     []string          `json:"include-repo,omitempty" yaml:"include-repo,omitempty" toml:"include-repo,omitempty"`
	Branch          map[string]string `json:"branch,omitempty" yaml:"branch,omitempty" toml:"branch,omitempty"`
	Tag             map[string]string `json:"tag,omitempty" yaml:"tag,omitempty" toml:"tag,omitempty"`
	ExcludeBranches []string          `json:"exclude-branches,omitempty" yaml:"exclude-branches,omitempty" toml:"exclude-branches,omitempty"`
}

func (p *CfgLazyGoMod) apply() {
	if len(p.Branch) == 0 {
		p.Branch = map[string]string{
			"master":  "master",
			"develop": "develop",
			"main":    "main",
			"release": "release",
			"test":    "test",
			"canary":  "canary",
			"dev":     "dev",
			"staging": "staging",
			"alpha":   "alpha",
			"beta":    "beta",
			"rc":      "rc",
		}
	}

	if len(p.Tag) == 0 {
		p.Tag = make(map[string]string)
	}

	if len(p.ExcludeBranches) == 0 {
		p.ExcludeBranches = make([]string, 0)
	}
}

type CfgLazyGen struct {
	Config *bool `json:"config,omitempty" yaml:"config,omitempty" toml:"config,omitempty"`
	Table  *bool `json:"table,omitempty" yaml:"table,omitempty" toml:"table,omitempty"`
	Cache  *bool `json:"cache,omitempty" yaml:"cache,omitempty" toml:"cache,omitempty"`
	I18n   *bool `json:"i18n,omitempty" yaml:"i18n,omitempty" toml:"i18n,omitempty"`
}

type CfgLazyConfig struct {
	GoMod *CfgLazyGoMod `json:"go-mod,omitempty" yaml:"go-mod,omitempty" toml:"go-mod,omitempty"`

	Gen *CfgLazyGen `json:"gen,omitempty" yaml:"gen,omitempty" toml:"gen,omitempty"`
}

func (c *CfgLazyConfig) Apply() {
	if c.GoMod == nil {
		c.GoMod = new(CfgLazyGoMod)
	}
	c.GoMod.apply()

	if c.Gen == nil {
		c.Gen = new(CfgLazyGen)
	}
}

var LazyConfig = new(CfgLazyConfig)

func LoadLazeConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Errorf("err:%v", err)
			return err
		}
		pterm.Warning.Printfln("%s is not found, use default", path)
	} else {
		defer file.Close()
		err = supportedExt["yaml"](file, &LazyConfig)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	LazyConfig.Apply()

	return nil
}
