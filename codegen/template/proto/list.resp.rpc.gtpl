message {{ .ResponseType }} {
	lazygophers.lrpc.core.Paginate paginate = 1;
	repeated {{ .Model }} list = 2;
}
