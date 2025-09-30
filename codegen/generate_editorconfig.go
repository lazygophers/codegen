package codegen

func GenerateEditorconfig(pb *PbPackage) (err error) {
	return generateConfigFile(pb, FileGeneratorConfig{
		PathType:     PathTypeEditorconfig,
		TemplateType: TemplateTypeEditorconfig,
		DisplayName:  ".editorconfig",
		Args:         nil,
	})
}
