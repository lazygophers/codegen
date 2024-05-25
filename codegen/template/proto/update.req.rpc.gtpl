message {{ .RequestType }} {
{{with .PprimaryKey }}	// @validate: required
	{{ $.PprimaryKeyType}} {{ $.PprimaryKey }} = 1;{{end}}
}
