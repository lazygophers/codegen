package state

import (
	"embed"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc"
	"github.com/lazygophers/lrpc/middleware/i18n"
)

//go:embed localize/
var i18nFs embed.FS

func LoadI18n() (err error) {
	State.I18n = i18n.DefaultI18n

	err = State.I18n.LoadLocalizes(i18nFs)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func Localize(key string, args ...interface{}) string {
	return State.I18n.Localize(key, args...)
}

func LocalizeWithCtx(ctx *lrpc.Ctx, key string, args ...interface{}) string {
	return State.I18n.LocalizeWithLanguage(i18n.ParseLanguage(ctx.Header("Accept-Language")), key, args...)
}
