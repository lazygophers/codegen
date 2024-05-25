{{ with .Role }}
	// @role: {{ $.Role }} {{end}}{{ with .Url }}
	// @url: {{ $.Url }} {{end}}{{ with .Model }}
	// @model: {{ $.Model }} {{end}}{{ with .Action }}
	// @action: {{ $.Action }} {{end}}
	rpc {{ .RpcName }} ({{ .RequestType }}) returns ({{ .ResponseType }}) {
	};
