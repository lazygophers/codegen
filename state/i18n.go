package state

import (
	"embed"
	"github.com/Xuanwo/go-locale"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/i18n"
	"github.com/pterm/pterm"
	"golang.org/x/text/language"
	"os"
)

var (
	//go:embed localize/*.yaml
	i18nFs embed.FS
	I18n   = i18n.DefaultI18n
)

func LoadI18n() (err error) {
	// 加载本地语言
	var defaultLanguage string
	if Config.Language != "" {
		defaultLanguage = Config.Language
	}

	if defaultLanguage == "" {
		defaultLanguage = os.Getenv("LANG")
	}

	if defaultLanguage == "" {
		lang, err := locale.Detect()
		if err != nil {
			log.Errorf("err:%v", err)
		} else {
			defaultLanguage = lang.String()
		}
	}

	if defaultLanguage == "" {
		defaultLanguage = language.English.String()
	}

	defaultLanguage = i18n.ParseLanguage(defaultLanguage)

	pterm.Info.Printfln("use language %s", defaultLanguage)
	I18n.SetDefaultLang(defaultLanguage)

	err = I18n.LoadLocalizes(i18nFs)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func Localize(key string, args ...interface{}) string {
	return I18n.Localize(key, args...)
}

func LocalizeWithLanguage(lang string, key string, args ...interface{}) string {
	return I18n.LocalizeWithLang(lang, key, args...)
}
