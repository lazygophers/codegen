package codegen

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/osx"
	"io/fs"
	"os"
)

func initStateDirectory(pb *PbPackage) error {
	if osx.IsDir(GetPath(PathTypeState, pb)) {
		return nil
	}

	err := os.MkdirAll(GetPath(PathTypeState, pb), fs.ModePerm)
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	return nil
}
