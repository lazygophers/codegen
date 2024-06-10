func {{ .RpcName }}(ctx *lrpc.Ctx, req *{{ .PB.GoPackageName }}.{{ .RequestType }}) (*{{ .PB.GoPackageName }}.{{ .ResponseType }}, error) {
	var rsp {{ .PB.GoPackageName }}.{{ .ResponseType }}

	return &rsp, nil
}
