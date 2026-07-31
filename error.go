package openingrouter

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type errorResponse struct {
	Error errorData `json:"error"`
}

type errorData struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Code    int64  `json:"code"`

	// per-kind union decoded lazily by toError
	Metadata json.RawMessage `json:"metadata"`
}

type errorMetadata struct {
	// provider errors
	Raw          string `json:"raw"`
	ProviderName string `json:"provider_name"`
	IsBYOK       bool   `json:"is_byok"`

	// typed provider errors
	ErrorType    ErrorType `json:"error_type"`
	ProviderCode string    `json:"provider_code"`

	// moderation errors
	Reasons      []string `json:"reasons"`
	FlaggedInput string   `json:"flagged_input"`
	ModelSlug    string   `json:"model_slug"`

	// credit and quota errors
	LimitSource string `json:"limit_source"`
	RemedyHint  string `json:"remedy_hint"`

	PreviousErrors []SubError `json:"previous_errors"`
}

var (
	ErrInvalidRequest      = errors.New("invalid request")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInsufficientCredits = errors.New("insufficient credits")
	ErrForbidden           = errors.New("forbidden")
	ErrNotFound            = errors.New("not found")
	ErrContextLength       = errors.New("context length exceeded")
	ErrModerated           = errors.New("moderated")
	ErrRateLimited         = errors.New("rate limited")
	ErrProviderUnavailable = errors.New("provider unavailable")
	ErrTimeout             = errors.New("timeout")
	ErrServer              = errors.New("server error")
)

func (e errorData) toError(status int, header http.Header) error {
	var meta errorMetadata

	if len(e.Metadata) != 0 {
		json.Unmarshal(e.Metadata, &meta)
	}

	code := e.Code
	if code == 0 {
		code = int64(status)
	}

	es := ErrorStatus{
		Code:         code,
		Type:         meta.ErrorType,
		ProviderCode: meta.ProviderCode,
		RetryAfter:   retryAfter(header),
		Previous:     meta.PreviousErrors,
	}

	for i := range es.Previous {
		previous := &es.Previous[i]

		if previous.Raw != "" {
			previous.Message = parseSubErrorMessage(previous.Raw)

			continue
		}

		previous.Message = parseSubErrorMessage(previous.Message)
	}

	message := parseSubErrorMessage(e.Message)

	switch {
	case meta.Raw != "":
		return &ProviderError{
			ErrorStatus:  es,
			Raw:          meta.Raw,
			ProviderName: meta.ProviderName,
			IsBYOK:       meta.IsBYOK,
		}
	case e.Name != "":
		return &ApiError{
			ErrorStatus: es,
			Name:        e.Name,
			Message:     message,
		}
	case len(meta.Reasons) != 0 || meta.FlaggedInput != "":
		if es.Type == "" {
			es.Type = ErrorTypeContentPolicy
		}

		return &ModerationError{
			ErrorStatus:  es,
			Message:      message,
			Reasons:      meta.Reasons,
			FlaggedInput: meta.FlaggedInput,
			ProviderName: meta.ProviderName,
			ModelSlug:    meta.ModelSlug,
		}
	case meta.LimitSource != "" || meta.RemedyHint != "" || code == http.StatusPaymentRequired:
		if es.Type == "" {
			es.Type = ErrorTypePaymentRequired
		}

		return &CreditsError{
			ErrorStatus: es,
			Message:     message,
			LimitSource: meta.LimitSource,
			RemedyHint:  meta.RemedyHint,
		}
	}

	return &OpenRouterError{
		ErrorStatus: es,
		Message:     message,
		Metadata:    e.Metadata,
	}
}

// AsOpenRouterError converts an HTTP response status or error into a structured OpenRouter error.
func AsOpenRouterError(resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var meta errorResponse

	err = json.NewDecoder(resp.Body).Decode(&meta)
	if err != nil {
		return newHttpError(resp)
	}

	info := meta.Error

	if info.Message == "" && info.Name == "" && len(info.Metadata) == 0 {
		return newHttpError(resp)
	}

	return info.toError(resp.StatusCode, resp.Header)
}

func newHttpError(resp *http.Response) error {
	return &HttpError{
		ErrorStatus: ErrorStatus{
			Code:       int64(resp.StatusCode),
			RetryAfter: retryAfter(resp.Header),
		},
		Status: resp.Status,
	}
}

func retryAfter(header http.Header) time.Duration {
	value := header.Get("Retry-After")
	if value == "" {
		return 0
	}

	seconds, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return max(time.Duration(seconds)*time.Second, 0)
	}

	date, err := http.ParseTime(value)
	if err != nil {
		return 0
	}

	return max(time.Until(date), 0)
}

func parseSubErrorMessage(raw string) string {
	if !strings.HasPrefix(raw, "{") && !strings.HasPrefix(raw, "[") {
		return raw
	}

	var root any

	err := json.Unmarshal([]byte(raw), &root)
	if err != nil {
		return raw
	}

	// direct string error message
	str, ok := root.(string)
	if ok {
		return str
	}

	// error object
	data, ok := root.(map[string]any)
	if !ok {
		// slice of error objects
		sl, ok := root.([]any)
		if !ok || len(sl) == 0 {
			return raw
		}

		data, ok = sl[0].(map[string]any)
		if !ok {
			return raw
		}
	}

	// sub-error object
	errorS, ok := data["error"]
	if ok {
		sub, ok := errorS.(map[string]any)
		if ok {
			data = sub
		}
	}

	// message string
	message, ok := data["message"]
	if ok {
		str, ok = message.(string)
		if ok {
			return str
		}
	}

	// error string
	errorS, ok = data["error"]
	if ok {
		str, ok = errorS.(string)
		if ok {
			return str
		}
	}

	return raw
}
