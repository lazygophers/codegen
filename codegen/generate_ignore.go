package codegen

func GenerateGitignote(pb *PbPackage) (err error) {
	return generateConfigFile(pb, FileGeneratorConfig{
		PathType:     PathTypeGitignore,
		TemplateType: TemplateTypeGitignore,
		DisplayName:  ".gitignore",
		Args:         nil,
	})
}

func GenerateDockerignote(pb *PbPackage) (err error) {
	return generateConfigFile(pb, FileGeneratorConfig{
		PathType:     PathTypeDockerignore,
		TemplateType: TemplateTypeDockerignore,
		DisplayName:  ".dockerignore",
		Args:         nil,
	})
}
