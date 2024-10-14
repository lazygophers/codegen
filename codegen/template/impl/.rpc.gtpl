func {{ .RpcName }}(ctx *lrpc.Ctx, req *{{ .RequestPackage }}.{{ .RequestType }}) (*{{ .ResponsePackage }}.{{ .ResponseType }}, error) {
	var rsp {{ .ResponsePackage }}.{{ .ResponseType }}

	return &rsp, nil
}
