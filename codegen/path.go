package codegen

import "path/filepath"

type PathType uint8

const (
	PathTypeRoot PathType = iota + 1
	PathTypePbGo
	PathTypeGoMod
	PathTypeConf

	PathTypeOrm
	PathTypeTableName
	PathTypeTableField

	PathTypeInternal

	PathTypeState
	PathTypeStateTable
	PathTypeStateConf
	PathTypeStateCache
	PathTypeStateState
	PathTypeStateI18n

	PathTypeImpl
	PathTypeImplPath
	PathTypeImplRoute

	PathTypeCmd
	PathTypeCmdMain

	PathTypeGoreleaser
	PathTypeMakefile
	PathTypeEditorconfig
	PathTypeGolangci
	PathTypeGitignore
	PathTypeDockerignore
)

func GetPath(t PathType, pb *PbPackage) string {
	switch t {
	case PathTypeRoot:
		return pb.ProjectRoot()

	case PathTypePbGo:
		return filepath.Join(pb.ProjectRoot(), pb.PackageName+".pb.go")

	case PathTypeGoMod:
		return filepath.Join(pb.ProjectRoot(), "go.mod")

	case PathTypeConf:
		return filepath.Join(pb.ProjectRoot(), "example.config.yaml")

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

	case PathTypeStateI18n:
		return filepath.Join(GetPath(PathTypeState, pb), "i18n.go")

	case PathTypeImpl:
		return filepath.Join(GetPath(PathTypeInternal, pb), "impl")

	case PathTypeImplPath:
		return filepath.Join(pb.ProjectRoot(), "rpc_path.gen.go")

	case PathTypeImplRoute:
		return filepath.Join(GetPath(PathTypeCmd, pb), "route.gen.go")

	case PathTypeCmd:
		return filepath.Join(pb.ProjectRoot(), "cmd")

	case PathTypeCmdMain:
		return filepath.Join(GetPath(PathTypeCmd, pb), "main.go")

	case PathTypeGoreleaser:
		return filepath.Join(pb.ProjectRoot(), ".goreleaser.yaml")

	case PathTypeMakefile:
		return filepath.Join(pb.ProjectRoot(), "Makefile")

	case PathTypeGolangci:
		return filepath.Join(pb.ProjectRoot(), ".golangci.yml")

	case PathTypeGitignore:
		return filepath.Join(pb.ProjectRoot(), ".gitignore")

	case PathTypeDockerignore:
		return filepath.Join(pb.ProjectRoot(), ".dockerignore")

	default:
		panic("unsupported path type")
	}
}
