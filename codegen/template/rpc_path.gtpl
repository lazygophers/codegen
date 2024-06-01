package {{ .PB.GoPackageName }}

const ({{ range $key, $value := .RPCS }}
    {{ $value.RpcName }} = "{{ $value.Path }}"{{ end }}
)
