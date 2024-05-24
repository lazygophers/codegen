package state

import "github.com/lazygophers/log"

type state struct {
	Config *Config
}

var State = new(state)

func Load() (err error) {
	err = LoadConfig()
	if err != nil {
	    log.Errorf("err:%v", err)
	    return err
	}

	err = ConnectCache()
	if err != nil {
	    log.Errorf("err:%v", err)
	    return err
	}

	err = ConnectCache()
	if err != nil {
	    log.Errorf("err:%v", err)
	    return err
	}

	return nil
}
