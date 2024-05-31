func {{ .RpcName }}(ctx *lrpc.Ctx, req *{{ .PB.GoPackageName }}.{{ .RequestType }}) (*{{ .PB.GoPackageName }}.{{ .ResponseType }}, error) {
	var rsp {{ .PB.GoPackageName }}.{{ .ResponseType }}

	err := state.{{ ToCamel (TrimPrefix .Model "Model") }}.
		NewScoop().
		Where("{{ .PrimaryKey }} = ?", req.{{ ToCamel .PrimaryKey }}).
		Delete().
		Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
