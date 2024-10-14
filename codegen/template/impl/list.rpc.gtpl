func {{ .RpcName }}(ctx *lrpc.Ctx, req *{{ .RequestPackage }}.{{ .RequestType }}) (*{{ .ResponsePackage }}.{{ .ResponseType }}, error) {
	var rsp {{ .ResponsePackage }}.{{ .ResponseType }}

	scoop := state.{{ ToCamel (TrimPrefix .Model "Model") }}.NewScoop()

	err := req.ListOption.Processor().{{ range $key, $value	:= .ListOption }}{{ if eq $value.FieldType "string" }}
		String(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value string) error {
			scoop.Where({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "string-slice" }}
		StringSlice(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value string) error {
			scoop.In({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "like" }}
		String(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value string) error {
			scoop.Like({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "left-like" }}
		String(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value string) error {
			scoop.LeftLike({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "right-like" }}
		String(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value string) error {
			scoop.RightLike({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "int" }}
		Int(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value int) error {
			scoop.Where({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "int8" }}
		Int8(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value int8) error {
			scoop.Where({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "int16" }}
		Int16(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value int16) error {
			scoop.Where({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "int32" }}
		Int32(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value int32) error {
			scoop.Where({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "int64" }}
		Int64(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value int64) error {
			scoop.Where({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "uint" }}
		Uint(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value uint) error {
			scoop.Where({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "uint8" }}
		Uint8(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value uint8) error {
			scoop.Where({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "uint16" }}
		Uint16(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value uint16) error {
			scoop.Where({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "uint32" }}
		Uint32(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value uint32) error {
			scoop.Where({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "uint64" }}
		Uint64(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value uint64) error {
			scoop.Where({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "float32" }}
		Float32(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value float32) error {
			scoop.Where({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "float64" }}
		Float64(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value float64) error {
			scoop.Where({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "bool" }}
		Bool(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value bool) error {
			scoop.Where({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "int-slice" }}
		IntSlice(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value []int) error {
			scoop.In({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "int8-slice" }}
		Int8Slice(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value []int8) error {
			scoop.In({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "int16-slice" }}
		Int16Slice(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value []int16) error {
			scoop.In({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "int32-slice" }}
		Int32Slice(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value []int32) error {
			scoop.In({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "int64-slice" }}
		Int64Slice(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value []int64) error {
			scoop.In({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "uint-slice" }}
		UintSlice(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value []uint) error {
			scoop.In({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "uint16-slice" }}
		Uint16Slice(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value []uint16) error {
			scoop.In({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "uint32-slice" }}
		Uint32Slice(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value []uint32) error {
			scoop.In({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "uint64-slice" }}
		Uint64Slice(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value []uint64) error {
			scoop.In({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "float32-slice" }}
		Float32Slice(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value []float32) error {
			scoop.In({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "float64-slice" }}
		Float64Slice(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value []float64) error {
			scoop.In({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ else if eq $value.FieldType "bool-slice" }}
		BoolSlice(int32({{ $.PB.GoPackageName }}.{{ $value.KeyField }}), func(value []bool) error {
			scoop.In({{ $.PB.GoPackageName }}.Db{{ $value.FieldName }}, value)
			return nil
		}).{{ end }}{{ end }}
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
