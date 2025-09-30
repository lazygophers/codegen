package i18n

import (
	"golang.org/x/text/language"
	"os"
	"resty.dev/v3"
	"sync"
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

// TranslateType represents the type of translator to use
type TranslateType string

const (
	// TranslateTypeGoogleFree represents Google's free translation service
	TranslateTypeGoogleFree TranslateType = "google-free"
)

// Translator is the interface that wraps the basic Translate method.
// Implementations should translate text from source language to destination language.
type Translator interface {
	Translate(srcLang, dstLang language.Tag, key string) (string, error)
}

var _ Translator = (*TranslatorHandler)(nil)

// TranslatorHandler wraps a handler function to implement Translator interface
type TranslatorHandler struct {
	handler func(srcLang, dstLang language.Tag, key string) (string, error)
}

var (
	translatorMap   = make(map[TranslateType]Translator)
	translatorMutex sync.RWMutex
	translatorOnce  = make(map[TranslateType]*sync.Once)
)

// RegisterTranslator registers a translator for the given type.
// It's safe to call from multiple goroutines.
func RegisterTranslator(tt TranslateType, t Translator) {
	translatorMutex.Lock()
	defer translatorMutex.Unlock()
	translatorMap[tt] = t
}

// GetTranslator returns the translator for the given type.
// It lazily initializes the translator on first access.
// It's safe to call from multiple goroutines.
func GetTranslator(tt TranslateType) (Translator, bool) {
	translatorMutex.RLock()
	t, ok := translatorMap[tt]
	translatorMutex.RUnlock()

	if ok {
		return t, true
	}

	// Lazy initialization for built-in translators
	translatorMutex.Lock()
	defer translatorMutex.Unlock()

	// Double-check after acquiring write lock
	if t, ok := translatorMap[tt]; ok {
		return t, true
	}

	// Initialize built-in translator if not exists
	if _, exists := translatorOnce[tt]; !exists {
		translatorOnce[tt] = &sync.Once{}
	}

	translatorOnce[tt].Do(func() {
		switch tt {
		case TranslateTypeGoogleFree:
			translatorMap[tt] = NewTranslatorGoogleFree()
		}
	})

	t, ok = translatorMap[tt]
	return t, ok
}

// Translate translates the given key from source language to destination language
func (p *TranslatorHandler) Translate(srcLang, dstLang language.Tag, key string) (string, error) {
	return p.handler(srcLang, dstLang, key)
}

// NewTranslatorHandler creates a new TranslatorHandler with the given handler function
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

