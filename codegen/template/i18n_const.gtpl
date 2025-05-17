package {{ .DirName }}

const ({{ range $value := .SortedKeys}}
	I18nTag{{ ToCamel $value }}		= `{{ $value }}`{{ end }}
)
