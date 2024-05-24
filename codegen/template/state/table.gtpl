package state

{{ if gt (len .GoImports) 0 }}import (
{{ range $key, $value := .GoImports }}    "{{ $value }}"
{{ end }}){{ end }}

var (
	_db *db.Client

{{ range $key, $value := .Models}}    {{TrimPrefix $value "Model"}} *db.Model[{{ $.PB.GoPackageName }}.{{ $value }}]
{{ end }})

func ConnectDatebase() (err error) {
	log.Info("try init database")
	_db, err = db.New(State.Config.Db,
	{{ range $key, $value := .Models}}    &{{ $.PB.GoPackageName }}.{{ $value }}{},
	{{ end }})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

{{ range $key, $value := .Models}}    {{TrimPrefix $value "Model"}} = db.NewModel[{{ $.PB.GoPackageName }}.{{ $value }}](Db()).SetNotFound(common.NewError({{ $.PB.GoPackageName }}.ErrCode_{{TrimPrefix $value "Model"}}NotFound))
{{ end }}
	log.Info("connect mysql successfully")

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
