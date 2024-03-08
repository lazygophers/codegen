package main

import (
	"github.com/lazygophers/codegen/cli"
	"github.com/lazygophers/utils/app"
)

func main() {
	if app.Name == "" {
		app.Name = "codegen"
	}

	cli.Run()
}
