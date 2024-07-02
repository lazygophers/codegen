package impl

import (
	"{{ .PB.GoPackage }}"
	"{{ .PB.GoPackage }}/internal/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc"
	"github.com/lazygophers/lrpc/middleware/xerror"
)
