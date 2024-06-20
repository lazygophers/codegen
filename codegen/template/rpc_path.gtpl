package {{ .PB.GoPackageName }}

const (
	ServerName = "{{ .PB.GoPackageName }}"
)

const ({{ range $key, $value := .RPCS }}
	RpcPath{{ $value.RpcName }} = "{{ $value.Path }}"{{ end }}
)
