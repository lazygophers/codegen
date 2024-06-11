package {{ .DirName }}

const ({{ range $key, $value := .Localize}}
	I18nTag{{ ToCamel $key }}		= `{{ $key }}`{{ end }}
)
