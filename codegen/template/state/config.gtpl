package state

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/cache"
	"github.com/lazygophers/utils/config"
	"github.com/lazygophers/utils/db"
)

type Config struct {
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
