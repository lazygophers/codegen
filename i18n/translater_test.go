package i18n_test

import (
	"github.com/lazygophers/codegen/i18n"
	"golang.org/x/text/language"
	"testing"
)

func TestDeelp(t *testing.T) {
	t.Log(i18n.NewTransacterDeeplFree().Translate(
		language.Chinese,
		language.English,
		"我们每一天都很懒惰",
	))
}
