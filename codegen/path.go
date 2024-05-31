package codegen

import "path/filepath"

type PathType uint8

const (
	PathTypePbGo PathType = iota + 1
	PathTypeGoMod
	PathTypeEditorconfig

	PathTypeOrm
	PathTypeTableName
	PathTypeTableField

	PathTypeInternal

	PathTypeState
	PathTypeStateTable
	PathTypeStateConf
	PathTypeStateCache
	PathTypeStateState

	PathTypeImpl
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
		return filepath.Join(pb.ProjectRoot(), "orm.gen.go")

	case PathTypeTableName:
		return filepath.Join(pb.ProjectRoot(), "table_name.gen.go")

	case PathTypeTableField:
		return filepath.Join(pb.ProjectRoot(), "table_field.gen.go")

	case PathTypeInternal:
		return filepath.Join(pb.ProjectRoot(), "internal")

	case PathTypeState:
		return filepath.Join(GetPath(PathTypeInternal, pb), "state")

	case PathTypeStateTable:
		return filepath.Join(GetPath(PathTypeState, pb), "table.go")

	case PathTypeStateConf:
		return filepath.Join(GetPath(PathTypeState, pb), "config.go")

	case PathTypeStateCache:
		return filepath.Join(GetPath(PathTypeState, pb), "cache.go")

	case PathTypeStateState:
		return filepath.Join(GetPath(PathTypeState, pb), "state.go")

	case PathTypeImpl:
		return filepath.Join(GetPath(PathTypeInternal, pb), "impl")

	default:
		panic("unsupported path type")
	}
}
