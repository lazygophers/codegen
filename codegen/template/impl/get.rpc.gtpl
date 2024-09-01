func {{ .RpcName }}(ctx *lrpc.Ctx, req *{{ .RequestPackage }}.{{ .RequestType }}) (*{{ .ResponsePackage }}.{{ .ResponseType }}, error) {
	var rsp {{ .ResponsePackage }}.{{ .ResponseType }}

	{{ ToSmallCamel (TrimPrefix .Model "Model") }}, err := state.{{ ToCamel (TrimPrefix .Model "Model") }}.
		NewScoop().
		Where("{{ .PrimaryKey }} = ?", req.{{ ToCamel .PrimaryKey }}).
		First()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.{{ ToCamel (TrimPrefix .Model "Model") }} = {{ ToSmallCamel (TrimPrefix .Model "Model") }}

	return &rsp, nil
}
