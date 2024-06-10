package state

import (
	"github.com/lazygophers/log"{{ with .Config.I18n }}
	"github.com/lazygophers/lrpc/middleware/i18n"{{end}}
)

type state struct {
	Config *Config

	// NOTE: Please fill in the state below
	{{ with .Config.I18n }}I18n *i18n.I18n{{ end }}
}

var State = new(state)

func Load() (err error) {
	err = LoadConfig()
	if err != nil {
	    log.Errorf("err:%v", err)
	    return err
	}{{ with .Config.I18n}}

	err = LoadI18n()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}{{end}}{{ with .Config.Cache }}

	err = ConnectCache()
	if err != nil {
	    log.Errorf("err:%v", err)
	    return err
	}{{ end }}{{ with .Config.Table}}

	err = ConnectDatebase()
	if err != nil {
	    log.Errorf("err:%v", err)
	    return err
	}{{ end }}

	return nil
}
