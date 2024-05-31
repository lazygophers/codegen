func {{ RpcName }}(ctx *lrpc.Ctx, req *{{ .PB.GoPackageName }}.{{ ReqType }}) (*{{ .PB.GoPackageName }}.{{ .RespType }}, error) {
	var resp {{ .PB.GoPackageName }}.{{ .RespType }}

	return &resp, nil
}
