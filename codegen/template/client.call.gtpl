func {{ .RPC.RpcName }}(ctx *lrpc.Ctx, req *{{ .RequestPackage }}.{{ .RequestType }}) (*{{ .ResponsePackage }}.{{ .ResponseType }}, error) {
	var rsp {{ .ResponsePackage }}.{{ .ResponseType }}
	return &rsp, lrpc.Call(ctx, &core.ServiceDiscoveryClient{
		ServiceName: ServerName,
		ServicePath: RpcPath{{ .RPC.RpcName }},
	}, req, &rsp)
}
