message {{ .RequestType }} {
{{with .PrimaryKey }}	// @validate: required
	{{ $.PrimaryKeyType}} {{ $.PrimaryKey }} = 1;{{end}}
}
