message {{ .ResponseType }} {
	{{ .Model }} {{ ToSnake (TrimPrefix .Model "Model") }} = 1;
}
