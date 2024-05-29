{{ with .Url }}
	// @url: {{ $.Url }} {{end}}{{ with .GenTo }}
	// @gen: {{ $.GenTo }} {{end}}{{ with .Model }}
	// @model: {{ $.Model }} {{end}}{{ with .Action }}
	// @action: {{ $.Action }} {{end}}
	rpc {{ .RpcName }} ({{ .RequestType }}) returns ({{ .ResponseType }}) {
		option (lrpc.lazygen) = { {{ with .Role }}
			role: "{{ $.Role }}",{{ end }}
		};
	};
