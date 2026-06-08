package state

import (
	"embed"
	"github.com/Xuanwo/go-locale"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/i18n"
	"github.com/lazygophers/utils/language"
	"github.com/pterm/pterm"
	"os"
)

var (
	//go:embed localize/*.yaml
	i18nFs embed.FS
	I18n   = i18n.Default
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
		defaultLanguage = "en"
	}

	defaultLanguage = language.Make(defaultLanguage).String()

	pterm.Info.Printfln("use language %s", defaultLanguage)
	I18n.SetDefaultLang(language.Make(defaultLanguage).Tag())

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
	return I18n.LocalizeWithLang(language.Make(lang).Tag(), key, args...)
}
