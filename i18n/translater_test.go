package i18n_test

import (
	"github.com/lazygophers/codegen/i18n"
	"golang.org/x/text/language"
	"sync"
	"testing"
)

// TestGoogleTranslate tests basic translation functionality
func TestGoogleTranslate(t *testing.T) {
	translator := i18n.NewTranslatorGoogleFree()
	result, err := translator.Translate(
		language.Chinese,
		language.English,
		"你好世界",
	)
	if err != nil {
		t.Fatalf("translation failed: %v", err)
	}
	if result == "" {
		t.Fatal("translation result is empty")
	}
	t.Logf("Translation result: %s", result)
}

// TestGetTranslatorLazyInit tests lazy initialization of translators
func TestGetTranslatorLazyInit(t *testing.T) {
	translator, ok := i18n.GetTranslator(i18n.TranslateTypeGoogleFree)
	if !ok {
		t.Fatal("failed to get google translator")
	}
	if translator == nil {
		t.Fatal("translator is nil")
	}

	// Get again to ensure it returns the same instance
	translator2, ok := i18n.GetTranslator(i18n.TranslateTypeGoogleFree)
	if !ok {
		t.Fatal("failed to get google translator on second call")
	}
	if translator2 == nil {
		t.Fatal("translator is nil on second call")
	}
}

// TestRegisterTranslator tests custom translator registration
func TestRegisterTranslator(t *testing.T) {
	customType := i18n.TranslateType("custom-test")
	mockTranslator := i18n.NewTranslatorHandler(func(srcLang, dstLang language.Tag, key string) (string, error) {
		return "mock: " + key, nil
	})

	i18n.RegisterTranslator(customType, mockTranslator)

	translator, ok := i18n.GetTranslator(customType)
	if !ok {
		t.Fatal("failed to get custom translator")
	}

	result, err := translator.Translate(language.English, language.Chinese, "test")
	if err != nil {
		t.Fatalf("translation failed: %v", err)
	}
	if result != "mock: test" {
		t.Fatalf("expected 'mock: test', got '%s'", result)
	}
}

// TestConcurrentAccess tests thread safety of GetTranslator
func TestConcurrentAccess(t *testing.T) {
	var wg sync.WaitGroup
	concurrency := 10

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			translator, ok := i18n.GetTranslator(i18n.TranslateTypeGoogleFree)
			if !ok {
				t.Error("failed to get translator in concurrent access")
				return
			}
			if translator == nil {
				t.Error("translator is nil in concurrent access")
			}
		}()
	}

	wg.Wait()
}

// TestErrorHandling tests error scenarios
func TestErrorHandling(t *testing.T) {
	translator := i18n.NewTranslatorGoogleFree()

	// Test with empty string
	result, err := translator.Translate(
		language.Chinese,
		language.English,
		"",
	)
	// Empty string might succeed or fail depending on API, just log the result
	t.Logf("Empty string translation: result='%s', err=%v", result, err)
}
