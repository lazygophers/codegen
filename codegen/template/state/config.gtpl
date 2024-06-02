package state

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/cache"
	"github.com/lazygophers/lrpc/middleware/db"
	"github.com/lazygophers/utils/config"
)

type Config struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`

	Port int `json:"port,omitempty" yaml:"port,omitempty" toml:"port,omitempty"`

	Db    *db.Config    `json:"db,omitempty" yaml:"db,omitempty" toml:"db,omitempty"`
	Cache *cache.Config `json:"cache,omitempty" yaml:"cache,omitempty" toml:"cache,omitempty"`

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
