message {{ .RequestType }} {
	// @validate: required
	{{ .Model }} {{ ToSnake (TrimPrefix .Model "Model") }} = 1;
}
