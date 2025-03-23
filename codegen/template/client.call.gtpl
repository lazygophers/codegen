func {{ .RPC.RpcName }}(ctx *lrpc.Ctx, req *{{ if ne .RequestPackage .PB.PackageName }}{{ .RequestPackage }}.{{ end }}{{ .RequestType }}) (*{{ if ne .ResponsePackage .PB.PackageName }}{{ .ResponsePackage }}.{{ end }}{{ .ResponseType }}, error) {
	var rsp {{ if ne .ResponsePackage .PB.PackageName }}{{ .ResponsePackage }}.{{ end }}{{ .ResponseType }}
	return &rsp, lrpc.Call(ctx, &core.ServiceDiscoveryClient{
		ServiceName: ServerName,
		ServicePath: RpcPath{{ .RPC.RpcName }},
	}, req, &rsp)
}
