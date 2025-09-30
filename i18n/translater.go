package i18n

import (
	"golang.org/x/text/language"
	"os"
	"resty.dev/v3"
	"time"
)

const (
	// HTTP client configuration
	defaultTimeout       = 10 * time.Second
	defaultRetryWaitTime = 1 * time.Second
	defaultRetryCount    = 3

	// Rate limit retry wait times
	googleRateLimitWait = 30 * time.Second
)

type TranslateType string

const (
	TranslateTypeGoogleFree TranslateType = "google-free"
)

type Translator interface {
	Translate(srcLang, dstLang language.Tag, key string) (string, error)
}

var _ Translator = (*TranslatorHandler)(nil)

type TranslatorHandler struct {
	handler func(srcLang, dstLang language.Tag, key string) (string, error)
}

var translatorMap = map[TranslateType]Translator{
	TranslateTypeGoogleFree: NewTranslatorGoogleFree(),
}

func RegisterTranslator(tt TranslateType, t Translator) {
	translatorMap[tt] = t
}

func GetTranslator(tt TranslateType) (Translator, bool) {
	t, ok := translatorMap[tt]
	return t, ok
}

func (p *TranslatorHandler) Translate(srcLang, dstLang language.Tag, key string) (string, error) {
	return p.handler(srcLang, dstLang, key)
}

func NewTranslatorHandler(handler func(srcLang, dstLang language.Tag, key string) (string, error)) *TranslatorHandler {
	return &TranslatorHandler{
		handler: handler,
	}
}

// setupProxy configures proxy settings from environment variables
func setupProxy(client *resty.Client) {
	if proxy := os.Getenv("HTTPS_PROXY"); proxy != "" {
		client.SetProxy(proxy)
	} else if proxy := os.Getenv("ALL_PROXY"); proxy != "" {
		client.SetProxy(proxy)
	} else if proxy := os.Getenv("HTTP_PROXY"); proxy != "" {
		client.SetProxy(proxy)
	}
}

