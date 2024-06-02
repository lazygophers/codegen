package main

import (
	"github.com/lazygophers/{{ $.PB.GoPackageName }}"
	"github.com/lazygophers/{{ $.PB.GoPackageName }}/internal/impl"
	"github.com/lazygophers/lrpc"
)

var Routes = []*lrpc.Route{ {{ range $key, $value := .RPCS }}
	{
		Method:  "{{ $value.Method }}",
		Path:    {{ $.PB.GoPackageName }}.RpcPath{{ $value.RpcName }},
		Handler: impl.{{ $value.RpcName }},
		Extra: map[string]any{
			"role": "{{ $value.Role }}",
		},
	},{{ end }}
}
