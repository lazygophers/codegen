package state

import (
	"github.com/lazygophers/log"{{ with .Config.Cache }}
	"github.com/lazygophers/lrpc/middleware/storage/cache"{{ end }}{{ with .Config.Table }}
	"github.com/lazygophers/lrpc/middleware/storage/db"{{ end }}
	"github.com/lazygophers/utils/config"
)

type Config struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`

	Port int `json:"port,omitempty" yaml:"port,omitempty" toml:"port,omitempty"`
{{ with .Config.Table }}	Db    *db.Config    `json:"db,omitempty" yaml:"db,omitempty" toml:"db,omitempty"`{{ end }}
{{ with .Config.Cache }}	Cache *cache.Config `json:"cache,omitempty" yaml:"cache,omitempty" toml:"cache,omitempty"`{{ end }}

	// NOTE: Please fill in the configuration below
}

func LoadConfig() (err error) {
	err = config.LoadConfig(State.Config)
		if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
