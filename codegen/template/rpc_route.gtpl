package main

import (
	"github.com/lazygophers/{{ $.PB.GoPackageName }}"
	"github.com/lazygophers/{{ $.PB.GoPackageName }}/internal/impl"
	"github.com/lazygophers/lrpc"
)

var app = lrpc.NewApp(&lrpc.Config{})

var Routes = []*lrpc.Route{ {{ range $key, $value := .RPCS }}
	{
		Method:  "{{ $value.Method }}",
		Path:    {{ $.PB.GoPackageName }}.RpcPath{{ $value.RpcName }},
		Handler: app.ToHandlerFunc(impl.{{ $value.RpcName }}),
		Extra: map[string]any{
			"role": "{{ $value.Role }}",
		},
	},{{ end }}
}
