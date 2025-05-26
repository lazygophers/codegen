package {{ .DirName }}

const ({{ range $key, $value := (Sort (MapKeysString .Localize))}}
	I18nTag{{ ToCamel $value }}		= `{{ $value }}`{{ end }}
)
