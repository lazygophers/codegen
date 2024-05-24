package {{ .PB.GoPackageName }}

{{ range $key, $value := .Models }}func ({{ $value }}) TableName() string {
	return "{{ $.PB.GoPackageName }}_{{ ToSnake (TrimPrefix $value "Model") }}"
}
{{ end }}
