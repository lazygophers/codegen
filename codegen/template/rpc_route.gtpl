package main

import ({{ if gt (len .RPCS) 0 }}
	"github.com/lazygophers/{{ $.PB.GoPackageName }}"
	"github.com/lazygophers/{{ $.PB.GoPackageName }}/internal/impl" {{ end }}
	"github.com/lazygophers/lrpc"
)

var app = lrpc.NewApp(&lrpc.Config{})

var Routes = []*lrpc.Route{ {{ range $key, $value := .RPCS }}{{ if not $value.SkipGenRoute }}
	{
		Method:  "{{ $value.Method }}",
		Path:    {{ $.PB.GoPackageName }}.RpcPath{{ $value.RpcName }},
		Handler: app.ToHandlerFunc(impl.{{ $value.RpcName }}),
		Extra: map[string]any{
			"role": "{{ $value.Role }}",
		},
	},{{ end }}{{ end }}
}
