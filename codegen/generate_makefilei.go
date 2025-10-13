package codegen

func GenerateMakefile(pb *PbPackage) (err error) {
	return generateConfigFile(pb, FileGeneratorConfig{
		PathType:     PathTypeMakefile,
		TemplateType: TemplateTypeMakefile,
		DisplayName:  "Makefile",
		Args:         nil,
	})
}
