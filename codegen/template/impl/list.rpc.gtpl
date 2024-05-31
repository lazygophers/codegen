func {{ .RpcName }}(ctx *lrpc.Ctx, req *{{ .PB.GoPackageName }}.{{ .RequestType }}) (*{{ .PB.GoPackageName }}.{{ .ResponseType }}, error) {
	var rsp {{ .PB.GoPackageName }}.{{ .ResponseType }}

	scoop := state.{{ ToCamel (TrimPrefix .Model "Model") }}.NewScoop()

	err := req.ListOption.Processor().
		//String(int32(example.ListTaskReq_ListOptionTitle), func(value string) error {
		//	scoop.Like(example.DbTitle, value)
		//	return nil
		//}).
		//String(int32(example.ListTaskReq_ListOptionUid), func(value string) error {
		//	scoop.Like(example.DbUid, value)
		//	return nil
		//}).
		Process()
	if err != nil {
		if xerror.CheckCode(err, xerror.ErrNoData) {
			rsp.Paginate = req.ListOption.Paginate()
			return &rsp, nil
		}

		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Paginate, rsp.List, err = scoop.FindByPage(req.ListOption)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
