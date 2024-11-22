package state

import (
	"{{ .PB.GoPackage }}"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/lrpc/middleware/storage/db"
)

var (
	_db *db.Client

{{ range $key, $value := .Models }}    {{TrimPrefix $value "Model"}} *db.Model[{{ $.PB.GoPackageName }}.{{ $value }}]
{{ end }})

func ConnectDatabase() (err error) {
	log.Info("try init database")
	_db, err = db.New(State.Config.Db,
	{{ range $key, $value := .Models}}    &{{ $.PB.GoPackageName }}.{{ $value }}{},
	{{ end }})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

{{ range $key, $value := .Models}}    {{TrimPrefix $value "Model"}} = db.NewModel[{{ $.PB.GoPackageName }}.{{ $value }}](Db()){{ if $.EnableErrorNotFound }}.
		SetNotFound(xerror.NewError(int32({{ $.PB.GoPackageName }}.ErrCode_{{TrimPrefix $value "Model"}}NotFound))){{ end }}{{ if $.EnableErrorDuplicateKey }}.
		SetDuplicatedKeyError(xerror.NewError(int32({{ $.PB.GoPackageName }}.ErrCode_{{TrimPrefix $value "Model"}}DuplicateKey))){{ end }}
{{ end }}
	log.Info("connect database successfully")

	return nil
}

func Db() *db.Client {
	return _db
}

func NewScoop() *db.Scoop {
	return _db.NewScoop()
}

func Begin() *db.Scoop {
	return NewScoop().Begin()
}

func CommitOrRollback(logic func(tx *db.Scoop) error) error {
	return NewScoop().CommitOrRollback(NewScoop().Begin(), logic)
}
