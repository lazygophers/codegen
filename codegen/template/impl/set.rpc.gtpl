func {{ .RpcName }}(ctx *lrpc.Ctx, req *{{ .PB.GoPackageName }}.{{ .RequestType }}) (*{{ .PB.GoPackageName }}.{{ .ResponseType }}, error) {
	var rsp {{ .PB.GoPackageName }}.{{ .ResponseType }}

	{{ ToSmallCamel (TrimPrefix .Model "Model") }} := req.{{ ToCamel (TrimPrefix .Model "Model") }}

	err := state.{{ ToCamel (TrimPrefix .Model "Model") }}.
		NewScoop().
		Where("{{ .PrimaryKey }} = ?", {{ ToSmallCamel (TrimPrefix .Model "Model") }}.{{ ToCamel .PrimaryKey }}).
		Updates(&{{ ToSmallCamel (TrimPrefix .Model "Model") }}).
		Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.{{ ToCamel (TrimPrefix .Model "Model") }}, err = state.{{ ToCamel (TrimPrefix .Model "Model") }}.
		NewScoop().
		Where("{{ .PrimaryKey }} = ?", {{ ToSmallCamel (TrimPrefix .Model "Model") }}.{{ ToCamel .PrimaryKey }}).
		First()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
