message {{ .ResponseType }} {
    core.Paginate paginate = 1;
	repeated {{ .Model }} list = 2;
}
