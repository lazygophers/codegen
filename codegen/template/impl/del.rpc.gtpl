func {{ .RpcName }}(ctx *lrpc.Ctx, req *{{ .RequestPackage }}.{{ .RequestType }}) (*{{ .ResponsePackage }}.{{ .ResponseType }}, error) {
	var rsp {{ .ResponsePackage }}.{{ .ResponseType }}

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
