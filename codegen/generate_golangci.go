package codegen

func GenerateGolangci(pb *PbPackage) (err error) {
	return generateConfigFile(pb, FileGeneratorConfig{
		PathType:     PathTypeGolangci,
		TemplateType: TemplateTypeGolangci,
		DisplayName:  ".golangci.yml",
		Args:         nil,
	})
}
