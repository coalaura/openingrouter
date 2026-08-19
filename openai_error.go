package openingrouter

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// openaiErrorBody is the root of an OpenAI compatible error response.
type openaiErrorBody struct {
	Error openaiErrorData `json:"error"`
}

// openaiErrorData is the error object of an OpenAI compatible error response.
// Code is decoded lazily because providers encode it either as a string or a
// number.
type openaiErrorData struct {
	Message string          `json:"message"`
	Type    string          `json:"type"`
	Param   string          `json:"param"`
	Code    json.RawMessage `json:"code"`
}

// OpenAIError represents an error response returned by an OpenAI compatible API.
type OpenAIError struct {
	ErrorStatus

	Message string
	Type    string
	Param   string
	Code    string
}

// Error returns the formatted string representation of the OpenAI error.
func (o *OpenAIError) Error() string {
	var sb strings.Builder

	sb.Grow(12 + len(o.Message) + 2 + len(o.Type) + 2 + len(o.Code))

	sb.WriteString("openai error")

	if o.Type != "" {
		sb.WriteString(" [")
		sb.WriteString(o.Type)
		sb.WriteString("]")
	}

	if o.Code != "" {
		sb.WriteString(" code ")
		sb.WriteString(o.Code)
	}

	sb.WriteString(": ")
	sb.WriteString(o.Message)

	return sb.String()
}

// AsOpenAIError converts an HTTP response status or error into a structured
// OpenAI compatible error.
func AsOpenAIError(resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var body openaiErrorBody

	decodeErr := json.NewDecoder(resp.Body).Decode(&body)
	if decodeErr != nil || body.Error.Message == "" && body.Error.Type == "" && len(body.Error.Code) == 0 {
		return newHttpError(resp)
	}

	info := body.Error

	var code int64

	raw := strings.TrimSpace(string(info.Code))
	if raw != "" && raw != "null" {
		trimmed := strings.Trim(raw, `"`)

		if parsed, perr := strconv.ParseInt(trimmed, 10, 64); perr == nil {
			code = parsed
		} else {
			code = int64(resp.StatusCode)
		}
	}

	if code == 0 {
		code = int64(resp.StatusCode)
	}

	return &OpenAIError{
		ErrorStatus: ErrorStatus{
			Code:       code,
			Type:       openaiErrorType(info.Type),
			RetryAfter: retryAfter(resp.Header),
		},
		Message: info.Message,
		Type:    info.Type,
		Param:   info.Param,
		Code:    strings.Trim(string(info.Code), `"`),
	}
}

// openaiErrorType maps a raw OpenAI error type onto the shared ErrorType
// vocabulary so errors.Is works uniformly.
func openaiErrorType(raw string) ErrorType {
	switch strings.ToLower(raw) {
	case "invalid_request_error", "invalid_prompt":
		return ErrorTypeInvalidRequest
	case "authentication_error":
		return ErrorTypeAuthentication
	case "permission_error", "permission_denied":
		return ErrorTypePermissionDenied
	case "rate_limit_error", "rate_limit_exceeded":
		return ErrorTypeRateLimitExceeded
	case "insufficient_quota", "payment_required":
		return ErrorTypePaymentRequired
	case "not_found_error", "model_not_found":
		return ErrorTypeNotFound
	case "server_error", "internal_error":
		return ErrorTypeServer
	case "timeout", "request_timeout":
		return ErrorTypeTimeout
	case "content_policy_violation", "refusal":
		return ErrorTypeContentPolicy
	}

	return ""
}
