package i18n

import (
	"errors"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/fake"
	"github.com/pterm/pterm"
	"golang.org/x/text/language"
	"io"
	"net/http"
	"os"
	"resty.dev/v3"
	"strings"
	"time"
)

const (
	// HTTP client configuration
	defaultTimeout       = 10 * time.Second
	defaultRetryWaitTime = 1 * time.Second
	defaultRetryCount    = 3

	// Rate limit retry wait times
	googleRateLimitWait = 30 * time.Second
	deeplRateLimitWait  = 1 * time.Minute
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
	TranslateTypeDeeplFree:  NewTranslatorDeeplFree(),
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

type TranslatorGoogleFree struct {
	client *resty.Client
}

type googleSentences struct {
	Trans   string  `json:"trans"`
	Orig    string  `json:"orig"`
	Backend float64 `json:"backend"`
}

type googleTranslateResp struct {
	Sentences  []googleSentences `json:"sentences"`
	Src        string            `json:"src"`
	Confidence float64           `json:"confidence"`
	Spell      struct {
	} `json:"spell"`
	LdResult struct {
		Srclangs            []string  `json:"srclangs"`
		SrclangsConfidences []float64 `json:"srclangs_confidences"`
		ExtendedSrclangs    []string  `json:"extended_srclangs"`
	} `json:"ld_result"`
}

func (p *TranslatorGoogleFree) Translate(srcLang, dstLang language.Tag, key string) (string, error) {
	var rsp googleTranslateResp
	resp, err := p.client.R().
		SetQueryParams(map[string]string{
			"sl": srcLang.String(),
			"tl": dstLang.String(),
			"q":  key,
		}).
		SetResult(&rsp).
		Get("")
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	if resp.StatusCode() != 200 {
		log.Errorf("err:%v", err)
		return "", errors.New(resp.Status())
	}

	if len(rsp.Sentences) == 0 {
		return "", errors.New("not get translations")
	}

	trans := candy.FirstOr(candy.SortUsing(rsp.Sentences, func(a, b googleSentences) bool {
		return a.Backend > b.Backend
	}), googleSentences{})

	if trans.Trans == "" {
		return "", errors.New("not get translations")
	}

	return trans.Trans, nil
}

func NewTranslatorGoogleFree() *TranslatorGoogleFree {
	p := &TranslatorGoogleFree{
		client: resty.New().
			SetTimeout(defaultTimeout).
			SetRetryWaitTime(defaultRetryWaitTime).
			SetRetryCount(defaultRetryCount).
			SetHeaders(map[string]string{
				"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36",
				"Keep-Alive": "5",
			}).
			AddRequestMiddleware(func(client *resty.Client, request *resty.Request) error {
				request.SetCookie(&http.Cookie{})
				request.SetHeader("User-Agent", fake.RandomUserAgent())
				return nil
			}).
			SetQueryParams(map[string]string{
				"client": "gtx",
				"ie":     "UTF-8",
				"dt":     "t",
				"dj":     "1",
			}).
			AddRetryConditions(func(response *resty.Response, err error) bool {
				if err != nil {
					return true
				}

				switch response.StatusCode() {
				case 200:
					return false
				case 429:
					wait := time.Duration(response.Request.Attempt) * googleRateLimitWait
					pterm.Warning.Printfln("google request too quick, wait %s before retrying", wait)
					log.Warnf("google request too quick, wait %s before retrying", wait)
					time.Sleep(wait)
					return true

				default:
					return false
				}
			}).
			SetLogger(log.Clone().SetOutput(io.Discard)).
			SetBaseURL("https://translate.google.com/translate_a/single"),
	}

	setupProxy(p.client)
	return p
}

type TranslatorDeeplFree struct {
	client *resty.Client
}

type DeeplFreeJobSentences struct {
	Text   string `json:"text"`
	Id     int    `json:"id"`
	Prefix string `json:"prefix"`
}

type DeeplFreeJob struct {
	Kind               string                   `json:"kind"`
	Sentences          []*DeeplFreeJobSentences `json:"sentences"`
	RawEnContextBefore []string                 `json:"raw_en_context_before"`
	RawEnContextAfter  []string                 `json:"raw_en_context_after"`
	//PreferredNumBeams  int                      `json:"preferred_num_beams"`
	Quality string `json:"quality"`
}

type DeeplFreeLangPreference struct {
	Weight  map[string]any `json:"weight"`
	Default string         `json:"default"`
}

type DeeplFreeLang struct {
	TargetLang         string                   `json:"target_lang"`
	Preference         *DeeplFreeLangPreference `json:"preference"`
	SourceLangComputed string                   `json:"source_lang_computed"`
}

type DeeplFreeCommonJobParams struct {
	//RegionalVariant string `json:"regionalVariant"`
	Mode        string `json:"mode"`
	BrowserType int    `json:"browserType"`
	TextType    string `json:"textType"`
}

type DeeplFreeParams struct {
	Jobs            []*DeeplFreeJob           `json:"jobs"`
	Lang            *DeeplFreeLang            `json:"lang"`
	Priority        int                       `json:"priority"`
	CommonJobParams *DeeplFreeCommonJobParams `json:"commonJobParams"`
	//Timestamp       int64                     `json:"timestamp"`
}

type DeeplFreeReq struct {
	Jsonrpc string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  *DeeplFreeParams `json:"params"`
	//Id      int64            `json:"id"`
}

type DeeplFreeResultTranslationsSentences struct {
	Text string `json:"text"`
	Ids  []int  `json:"ids"`
}

type DeeplFreeResultTranslationsBeams struct {
	Sentences  []*DeeplFreeResultTranslationsSentences `json:"sentences"`
	NumSymbols int                                     `json:"num_symbols"`
}

type DeeplFreeResultTranslations struct {
	Beams   []*DeeplFreeResultTranslationsBeams `json:"beams"`
	Quality string                              `json:"quality"`
}

type DeeplFreeResult struct {
	Translations          []*DeeplFreeResultTranslations `json:"translations"`
	TargetLang            string                         `json:"target_lang"`
	SourceLang            string                         `json:"source_lang"`
	SourceLangIsConfident bool                           `json:"source_lang_is_confident"`
	DetectedLanguages     struct {
	} `json:"detectedLanguages"`
}

type DeeplFreeResp struct {
	Jsonrpc string           `json:"jsonrpc"`
	Result  *DeeplFreeResult `json:"result"`
}

func (p *TranslatorDeeplFree) Translate(srcLang, dstLang language.Tag, key string) (string, error) {
	var rsp DeeplFreeResp
	var resp, err = p.client.R().
		SetBody(&DeeplFreeReq{
			Jsonrpc: "2.0",
			Method:  "LMT_handle_jobs",
			Params: &DeeplFreeParams{
				Jobs: []*DeeplFreeJob{
					{
						Kind: "default",
						Sentences: []*DeeplFreeJobSentences{
							{
								Text:   key,
								Id:     1,
								Prefix: "",
							},
						},
						RawEnContextBefore: []string{},
						RawEnContextAfter:  []string{},
						Quality:            "fast",
					},
				},
				Lang: &DeeplFreeLang{
					TargetLang: strings.ToUpper(dstLang.String()),
					Preference: &DeeplFreeLangPreference{
						Weight:  map[string]any{},
						Default: "default",
					},
					SourceLangComputed: strings.ToUpper(srcLang.String()),
				},
				Priority: -1,
				CommonJobParams: &DeeplFreeCommonJobParams{
					//RegionalVariant: "",
					Mode:        "translate",
					BrowserType: 1,
					TextType:    "plaintext",
				},
				//Timestamp: time.Now().UnixMilli(),
			},
		}).
		SetResult(&rsp).
		Post("")
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	if resp.StatusCode() != 200 {
		log.Errorf("err:%v", err)
		return "", errors.New(resp.Status())
	}

	if rsp.Result == nil || len(rsp.Result.Translations) == 0 || len(rsp.Result.Translations[0].Beams) == 0 {
		return "", errors.New("not get translations")
	}

	trans := candy.FirstOr(candy.SortUsing(rsp.Result.Translations[0].Beams, func(a, b *DeeplFreeResultTranslationsBeams) bool {
		return a.NumSymbols > b.NumSymbols
	}), &DeeplFreeResultTranslationsBeams{})

	if len(trans.Sentences) == 0 || trans.Sentences[0].Text == "" {
		return "", errors.New("not get translations")
	}

	return trans.Sentences[0].Text, nil
}

func NewTranslatorDeeplFree() *TranslatorDeeplFree {
	p := &TranslatorDeeplFree{
		client: resty.New().
			SetTimeout(defaultTimeout).
			SetRetryWaitTime(defaultRetryWaitTime).
			SetRetryCount(defaultRetryCount).
			SetHeaders(map[string]string{
				"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.162 Safari/537.36",
				"Origin":     "https://www.deepl.com/",
				"Host":       "www2.deepl.com",
			}).
			AddRequestMiddleware(func(client *resty.Client, request *resty.Request) error {
				request.SetCookie(&http.Cookie{})
				request.SetHeader("User-Agent", fake.RandomUserAgent())
				return nil
			}).
			AddRetryConditions(func(response *resty.Response, err error) bool {
				if err != nil {
					return true
				}

				switch response.StatusCode() {
				case 200:
					return false
				case 429:
					log.Warnf("deepl request too quick, wait %s before retrying", deeplRateLimitWait)
					time.Sleep(deeplRateLimitWait)
					return true

				default:
					return false
				}
			}).
			SetQueryParams(map[string]string{
				"method": "LMT_handle_jobs",
			}).
			SetLogger(log.Clone().SetOutput(io.Discard)).
			SetBaseURL("https://www2.deepl.com/jsonrpc"),
	}

	setupProxy(p.client)
	return p
}
