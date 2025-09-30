package i18n

import (
	"fmt"
	"github.com/lazygophers/log"
	"golang.org/x/text/language"
	"strings"
)

// Language represents a language with its original string and parsed tag
type Language struct {
	Lang string
	Tag  language.Tag
}

// String returns the string representation of the language
func (l *Language) String() string {
	if l.Tag.String() == l.Lang {
		return l.Lang
	}

	return fmt.Sprintf("%s(%s)", l.Tag.String(), l.Lang)
}

// GoString implements fmt.GoStringer interface
func (l *Language) GoString() string {
	return l.String()
}

// ParseLanguage parses a language string into a Language struct.
// It supports various Chinese language variants and falls back to standard language parsing.
func ParseLanguage(lang string) (*Language, error) {
	lang = strings.ToLower(lang)

	// Map of language aliases to their standard tags
	var tag language.Tag
	switch lang {
	case "zh-cn", "zh-chs", "zh-hans":
		tag = language.Chinese
	case "zh-hk", "zh-tw", "zh-mo", "zh-sg", "zh-cht":
		tag = language.TraditionalChinese
	default:
		// Parse using standard library
		parsedTag, err := language.Parse(lang)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		tag = parsedTag
	}

	return &Language{
		Lang: lang,
		Tag:  tag,
	}, nil
}

// MustParseLanguage is like ParseLanguage but panics on error
func MustParseLanguage(lang string) *Language {
	l, err := ParseLanguage(lang)
	if err != nil {
		panic(err)
	}

	return l
}
