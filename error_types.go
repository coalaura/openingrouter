package openingrouter

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// SubError is one failed upstream attempt from metadata.previous_errors.
type SubError struct {
	Code         int64  `json:"code"`
	Message      string `json:"message"`
	ProviderName string `json:"provider_name"`
	Raw          string `json:"raw"`
}

// ErrorType is the stable error_type vocabulary OpenRouter tags errors with.
// It is more precise than the HTTP status and is carried identically on
// response bodies and mid-stream SSE events.
type ErrorType string

const (
	ErrorTypeContextLengthExceeded ErrorType = "context_length_exceeded"
	ErrorTypeMaxTokensExceeded     ErrorType = "max_tokens_exceeded"
	ErrorTypeTokenLimitExceeded    ErrorType = "token_limit_exceeded"
	ErrorTypeStringTooLong         ErrorType = "string_too_long"
	ErrorTypeAuthentication        ErrorType = "authentication"
	ErrorTypePermissionDenied      ErrorType = "permission_denied"
	ErrorTypePaymentRequired       ErrorType = "payment_required"
	ErrorTypeRateLimitExceeded     ErrorType = "rate_limit_exceeded"
	ErrorTypeProviderOverloaded    ErrorType = "provider_overloaded"
	ErrorTypeProviderUnavailable   ErrorType = "provider_unavailable"
	ErrorTypeInvalidRequest        ErrorType = "invalid_request"
	ErrorTypeInvalidPrompt         ErrorType = "invalid_prompt"
	ErrorTypeNotFound              ErrorType = "not_found"
	ErrorTypePreconditionFailed    ErrorType = "precondition_failed"
	ErrorTypePayloadTooLarge       ErrorType = "payload_too_large"
	ErrorTypeUnprocessable         ErrorType = "unprocessable"
	ErrorTypeContentPolicy         ErrorType = "content_policy_violation"
	ErrorTypeRefusal               ErrorType = "refusal"
	ErrorTypeInvalidImage          ErrorType = "invalid_image"
	ErrorTypeImageTooLarge         ErrorType = "image_too_large"
	ErrorTypeImageTooSmall         ErrorType = "image_too_small"
	ErrorTypeUnsupportedImage      ErrorType = "unsupported_image_format"
	ErrorTypeImageNotFound         ErrorType = "image_not_found"
	ErrorTypeImageDownloadFailed   ErrorType = "image_download_failed"
	ErrorTypeServer                ErrorType = "server"
	ErrorTypeTimeout               ErrorType = "timeout"
	ErrorTypeUnmapped              ErrorType = "unmapped"
)

// ErrorStatus carries the fields every OpenRouter error response shares. It is
// embedded in each concrete error type, so Code, Type, RetryAfter, Retryable
// and errors.Is work uniformly regardless of which one you got back.
type ErrorStatus struct {
	Code         int64
	Type         ErrorType
	ProviderCode string
	RetryAfter   time.Duration
	Previous     []SubError
}

// Unwrap maps the error onto its sentinel category. error_type wins when
// present because it survives provider-specific status mangling.
func (s ErrorStatus) Unwrap() error {
	switch s.Type {
	case ErrorTypeContextLengthExceeded, ErrorTypeMaxTokensExceeded, ErrorTypeStringTooLong:
		return ErrContextLength
	case ErrorTypeTokenLimitExceeded, ErrorTypePaymentRequired:
		return ErrInsufficientCredits
	case ErrorTypeAuthentication:
		return ErrUnauthorized
	case ErrorTypePermissionDenied:
		return ErrForbidden
	case ErrorTypeRateLimitExceeded:
		return ErrRateLimited
	case ErrorTypeProviderOverloaded, ErrorTypeProviderUnavailable:
		return ErrProviderUnavailable
	case ErrorTypeContentPolicy, ErrorTypeRefusal:
		return ErrModerated
	case ErrorTypeNotFound, ErrorTypeImageNotFound:
		return ErrNotFound
	case ErrorTypeTimeout:
		return ErrTimeout
	case ErrorTypeServer, ErrorTypeUnmapped:
		return ErrServer
	case ErrorTypeInvalidRequest, ErrorTypeInvalidPrompt, ErrorTypePreconditionFailed,
		ErrorTypePayloadTooLarge, ErrorTypeUnprocessable, ErrorTypeInvalidImage,
		ErrorTypeImageTooLarge, ErrorTypeImageTooSmall, ErrorTypeUnsupportedImage,
		ErrorTypeImageDownloadFailed:
		return ErrInvalidRequest
	}

	switch s.Code {
	case http.StatusBadRequest, http.StatusPreconditionFailed, http.StatusRequestEntityTooLarge, http.StatusUnprocessableEntity:
		return ErrInvalidRequest
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusPaymentRequired:
		return ErrInsufficientCredits
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusRequestTimeout, http.StatusGatewayTimeout:
		return ErrTimeout
	case http.StatusTooManyRequests:
		return ErrRateLimited
	case http.StatusBadGateway, http.StatusServiceUnavailable:
		return ErrProviderUnavailable
	case http.StatusInternalServerError:
		return ErrServer
	}

	return nil
}

// Retryable reports whether retrying the same request could plausibly succeed.
// Honor RetryAfter first when it is non-zero.
func (s ErrorStatus) Retryable() bool {
	switch s.Code {
	case http.StatusRequestTimeout, http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	}

	return false
}

// OpenRouterError represents an error response returned by the OpenRouter API.
type OpenRouterError struct {
	ErrorStatus

	Message string

	// Metadata holds metadata shapes this package does not model yet (e.g.
	// guardrail patterns), preserved verbatim rather than dropped.
	Metadata json.RawMessage
}

// Error returns the formatted string representation of the OpenRouter error.
func (o *OpenRouterError) Error() string {
	var (
		sb  strings.Builder
		buf [20]byte
	)

	sb.Grow(16 + 20 + 2 + len(o.Message))

	sb.WriteString("openrouter code ")
	sb.Write(strconv.AppendInt(buf[:0], o.Code, 10))
	sb.WriteString(": ")
	sb.WriteString(o.Message)

	return sb.String()
}

// ApiError represents a general API error returned by OpenRouter.
type ApiError struct {
	ErrorStatus

	Name    string
	Message string
}

// Error returns the formatted string representation of the API error.
func (a *ApiError) Error() string {
	var sb strings.Builder

	sb.Grow(10 + len(a.Name) + 2 + len(a.Message))

	sb.WriteString("api error ")
	sb.WriteString(a.Name)
	sb.WriteString(": ")
	sb.WriteString(a.Message)

	return sb.String()
}

// ProviderError represents an error returned by an underlying model provider.
type ProviderError struct {
	ErrorStatus

	Raw          string
	ProviderName string
	IsBYOK       bool
}

// Error returns the formatted string representation of the provider error.
func (p *ProviderError) Error() string {
	var sb strings.Builder

	message := parseSubErrorMessage(p.Raw)

	sb.Grow(17 + len(message))

	sb.WriteString("provider error: ")
	sb.WriteString(message)

	return sb.String()
}

// CreditsError represents a 402 caused by an insufficient credit balance.
type CreditsError struct {
	ErrorStatus

	Message     string
	LimitSource string
	RemedyHint  string
}

// Error returns the formatted string representation of the credits error.
func (c *CreditsError) Error() string {
	var sb strings.Builder

	sb.Grow(22 + len(c.Message))

	sb.WriteString("insufficient credits: ")
	sb.WriteString(c.Message)

	return sb.String()
}

// ModerationError represents input that was flagged before reaching a provider.
type ModerationError struct {
	ErrorStatus

	Message      string
	Reasons      []string
	FlaggedInput string
	ProviderName string
	ModelSlug    string
}

// Error returns the formatted string representation of the moderation error.
func (m *ModerationError) Error() string {
	var sb strings.Builder

	size := 18 + len(m.Message)

	for _, reason := range m.Reasons {
		size += len(reason) + 2
	}

	sb.Grow(size + 2)

	sb.WriteString("moderation error: ")
	sb.WriteString(m.Message)

	if len(m.Reasons) != 0 {
		sb.WriteString(" (")

		for i, reason := range m.Reasons {
			if i != 0 {
				sb.WriteString(", ")
			}

			sb.WriteString(reason)
		}

		sb.WriteString(")")
	}

	return sb.String()
}

// HttpError represents a response that carried no usable OpenRouter error body.
type HttpError struct {
	ErrorStatus

	Status string
}

// Error returns the formatted string representation of the HTTP error.
func (h *HttpError) Error() string {
	var sb strings.Builder

	sb.Grow(6 + len(h.Status))

	sb.WriteString("http: ")
	sb.WriteString(h.Status)

	return sb.String()
}
