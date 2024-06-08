package state

import "github.com/lazygophers/log"

func Load() (err error) {
	err = LoadConfig()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = LoadI18n()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
