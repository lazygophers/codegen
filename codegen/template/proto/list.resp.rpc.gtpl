message {{ .ResponseType }} {
	lrpc.Paginate paginate = 1;
	{{ .Model }} {{ ToSnake (TrimPrefix .Model "Model") }} = 2;
}
