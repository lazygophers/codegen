{{ with .Url }}
	// @url: {{ $.Url }} {{end}}{{ with .GenTo }}
	// @gen: {{ $.GenTo }} {{end}}{{ with .Model }}
	// @model: {{ $.Model }} {{end}}{{ with .Action }}
	// @action: {{ $.Action }} {{end}}
	rpc {{ .RpcName }} ({{ .RequestType }}) returns ({{ .ResponseType }}) {
			option (core.http) = { {{ with .Method }}
				method: "{{ $.Method }}",{{ end }}{{ with Path }}
				path: "{{ $.Path }}",{{ end }}
			};
		option (core.lazygen) = { {{ with .Role }}
			role: "{{ $.Role }}",{{ end }}
		};
	};
