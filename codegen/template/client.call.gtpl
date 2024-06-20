func {{ .RPC.RpcName }}(ctx *lrpc.Ctx, req *{{ .RequestType }}) (*{{ .ResponseType }}, error) {
	var rsp {{ .ResponseType }}
	return &rsp, lrpc.Call(ctx, &core.ServiceDiscoveryClient{
		ServiceName: ServerName,
		ServicePath: RpcPath{{ .RPC.RpcName }},
	}, req, &rsp)
}
