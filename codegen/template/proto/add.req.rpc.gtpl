message {{ .RequestType }} {
	{{ .Model }} {{ ToSnake (TrimPrefix .Model "Model") }} = 1;
}
