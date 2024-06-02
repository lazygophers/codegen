package main

import (
	"github.com/lazygophers/{{ .PB.GoPackageName }}"
	"github.com/lazygophers/{{ .PB.GoPackageName }}/internal/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc"
)

func main() {
	err := state.Load()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	app := lrpc.NewApp(&lrpc.Config{
		Name: state.State.Config.Name,
	})

	app.AddRoutes({{ .PB.GoPackageName }}.Routes)

	err = app.ListenAndServe(state.State.Config.Port)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
