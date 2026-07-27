package openingrouter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type errorResponse struct {
	Error errorData `json:"error"`
}

type errorData struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Code    int64  `json:"code"`

	Metadata *ProviderError `json:"metadata"`
}

type OpenRouterError struct {
	Message string
	Code    int64
}

type ApiError struct {
	Name    string
	Message string
}

type ProviderError struct {
	Raw          string `json:"raw"`
	ProviderName string `json:"provider_name"`
	IsBYOK       bool   `json:"is_byok"`
}

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

func (a *ApiError) Error() string {
	var sb strings.Builder

	sb.Grow(10 + len(a.Name) + 2 + len(a.Message))

	sb.WriteString("api error ")
	sb.WriteString(a.Name)
	sb.WriteString(": ")
	sb.WriteString(a.Message)

	return sb.String()
}

func (p *ProviderError) Error() string {
	var sb strings.Builder

	message := parseSubErrorMessage(p.Raw)

	sb.Grow(17 + len(message))

	sb.WriteString("provider error: ")
	sb.WriteString(message)

	return sb.String()
}

func AsOpenRouterError(resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var meta errorResponse

	err = json.NewDecoder(resp.Body).Decode(&meta)
	if err != nil {
		return fmt.Errorf("http: %s", resp.Status)
	}

	info := meta.Error

	if info.Metadata != nil {
		return info.Metadata
	}

	if info.Name != "" {
		return &ApiError{
			Name:    info.Name,
			Message: parseSubErrorMessage(info.Message),
		}
	}

	return &OpenRouterError{
		Message: parseSubErrorMessage(info.Message),
		Code:    info.Code,
	}
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
