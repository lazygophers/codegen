package {{ .PB.GoPackageName }}

const ({{range $key, $value := .FieldMap }}
	Db{{ ToCamel $key }} = "{{ $key }}"{{ end }}
)
