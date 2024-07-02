{{ with .Url }}
	// @url: {{ $.Url }} {{end}}{{ with .GenTo }}
	// @gen: {{ $.GenTo }} {{end}}{{ with .Model }}
	// @model: {{ $.Model }} {{end}}{{ with .Action }}
	// @action: {{ $.Action }} {{end}}
	rpc {{ .RpcName }} ({{ .RequestType }}) returns ({{ .ResponseType }}) {
			option (lazygophers.lrpc.core.http) = { {{ with .Method }}
				method: "{{ $.Method }}",{{ end }}{{ with .Path }}
				path: "{{ $.Path }}",{{ end }}
			};
		option (lazygophers.lrpc.core.lazygen) = { {{ with .Role }}
			role: "{{ $.Role }}",{{ end }}
		};
	};
