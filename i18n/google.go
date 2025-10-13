package i18n

import (
	"errors"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/candy"
	"github.com/lazygophers/utils/fake"
	"github.com/pterm/pterm"
	"golang.org/x/text/language"
	"io"
	"resty.dev/v3"
	"time"
)

// TranslatorGoogleFree implements Translator interface using Google's free translation service
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
	Spell      interface{}       `json:"spell"`
	LdResult   struct {
		Srclangs            []string  `json:"srclangs"`
		SrclangsConfidences []float64 `json:"srclangs_confidences"`
		ExtendedSrclangs    []string  `json:"extended_srclangs"`
	} `json:"ld_result"`
}

// Translate translates text from source language to destination language using Google Translate
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
		log.Errorf("http request failed with status %s", resp.Status())
		return "", errors.New(resp.Status())
	}

	if len(rsp.Sentences) == 0 {
		return "", errors.New("no translations found")
	}

	trans := candy.FirstOr(candy.SortUsing(rsp.Sentences, func(a, b googleSentences) bool {
		return a.Backend > b.Backend
	}), googleSentences{})

	if trans.Trans == "" {
		return "", errors.New("no translations found")
	}

	return trans.Trans, nil
}

// NewTranslatorGoogleFree creates a new Google translator instance
func NewTranslatorGoogleFree() *TranslatorGoogleFree {
	p := &TranslatorGoogleFree{
		client: resty.New().
			SetTimeout(defaultTimeout).
			SetRetryWaitTime(defaultRetryWaitTime).
			SetRetryCount(defaultRetryCount).
			SetHeaders(map[string]string{
				"Keep-Alive": "5",
			}).
			AddRequestMiddleware(func(client *resty.Client, request *resty.Request) error {
				// Set random User-Agent for each request to avoid rate limiting
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