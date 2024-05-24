package codegen

import "path/filepath"

type PathType uint8

const (
	PathTypePbGo PathType = iota + 1
	PathTypeGoMod
	PathTypeEditorconfig
	PathTypeOrm

	PathTypeState
	PathTypeStateTable
	PathTypeStateConf
	PathTypeStateCache
	PathTypeStateState
)

func GetPath(t PathType, pb *PbPackage) string {
	switch t {
	case PathTypePbGo:
		return filepath.Join(pb.ProjectRoot(), pb.PackageName+".pb.go")

	case PathTypeGoMod:
		return filepath.Join(pb.ProjectRoot(), "go.mod")

	case PathTypeEditorconfig:
		return filepath.Join(pb.ProjectRoot(), ".editorconfig")

	case PathTypeOrm:
		return filepath.Join(pb.ProjectRoot(), "orm_gen.go")

	case PathTypeState:
		return filepath.Join(pb.ProjectRoot(), "state")

	case PathTypeStateTable:
		return filepath.Join(pb.ProjectRoot(), "state", "table.go")

	case PathTypeStateConf:
		return filepath.Join(pb.ProjectRoot(), "state", "config.go")

	case PathTypeStateCache:
		return filepath.Join(pb.ProjectRoot(), "state", "cache.go")

	case PathTypeStateState:
		return filepath.Join(pb.ProjectRoot(), "state", "state.go")

	default:
		panic("unsupported path type")
	}
}
