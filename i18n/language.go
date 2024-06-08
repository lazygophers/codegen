package i18n

import (
	"fmt"
	"github.com/lazygophers/log"
	"golang.org/x/text/language"
	"strings"
)

type Language struct {
	Lang string
	Tag  language.Tag
}

func (l *Language) String() string {
	if l.Tag.String() == l.Lang {
		return l.Lang
	}

	return fmt.Sprintf("%s(%s)", l.Tag.String(), l.Lang)
}

func (l *Language) GoString() string {
	return l.String()
}

func (l *Language) ʔ() string {
	return l.String()
}

func ParseLanguage(lang string) (*Language, error) {
	lang = strings.ToLower(lang)

	switch lang {
	case "zh-cn", "zh-chs":
		return &Language{
			Lang: lang,
			Tag:  language.Chinese,
		}, nil

	case "zh-hk", "zh-tw", "zh-mo", "zh-sg", "zh-cht":
		return &Language{
			Lang: lang,
			Tag:  language.TraditionalChinese,
		}, nil

	default:
		l, err := language.Parse(lang)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		return &Language{
			Lang: lang,
			Tag:  l,
		}, nil
	}
}

func MustParseLanguage(lang string) *Language {
	l, err := ParseLanguage(lang)
	if err != nil {
		panic(err)
	}

	return l
}
