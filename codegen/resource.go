package codegen

import (
	"embed"
	"github.com/lazygophers/log"
)

//go:embed resource/*
var embedFs embed.FS

func GetDefaultEditorconfig() []byte {
	buf, err := embedFs.ReadFile("resource/.editorconfig")
	if err != nil {
		log.Panicf("err:%v", err)
		return nil
	}

	return buf
}
