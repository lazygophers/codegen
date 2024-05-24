package state

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/cache"
)

var (
	_cache cache.Cache
)

func ConnectCache() (err error) {
	log.Info("try init cache")
	_cache, err = cache.New(State.Config.Cache)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func Cache() cache.Cache {
	return _cache
}
