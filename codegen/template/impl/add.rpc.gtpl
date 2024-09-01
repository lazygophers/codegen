func {{ .RpcName }}(ctx *lrpc.Ctx, req *{{ .PB.GoPackageName }}.{{ .RequestType }}) (*{{ .PB.GoPackageName }}.{{ .ResponseType }}, error) {
	var rsp {{ .PB.GoPackageName }}.{{ .ResponseType }}

	{{ ToSmallCamel (TrimPrefix .Model "Model") }} := *req.{{ ToCamel (TrimPrefix .Model "Model") }}
	{{ ToSmallCamel (TrimPrefix .Model "Model") }}.{{ ToCamel .PrimaryKey }} = 0

	err := state.{{ ToCamel (TrimPrefix .Model "Model") }}.
		NewScoop().
		Create(&{{ ToSmallCamel (TrimPrefix .Model "Model") }})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.{{ ToCamel (TrimPrefix .Model "Model") }} = &{{ ToSmallCamel (TrimPrefix .Model "Model") }}

	return &rsp, nil
}
