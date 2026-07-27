package openingrouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type errorResponse[T any] struct {
	Error errorData[T] `json:"error"`
}

type errorData[T any] struct {
	Message  string `json:"message"`
	Code     int64  `json:"code"`
	Metadata T      `json:"metadata"`
}

type OpenRouterError struct {
	Message string
	Code    int64
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

func (p *ProviderError) Data() any {
	if p.Raw == "" {
		return nil
	}

	var data any

	err := json.Unmarshal([]byte(p.Raw), &data)
	if err != nil {
		return nil
	}

	return data
}

func (p *ProviderError) Error() string {
	var sb strings.Builder

	sb.Grow(len(p.ProviderName) + 17 + len(p.Raw))

	sb.WriteString(p.ProviderName)
	sb.WriteString(" provider error: ")
	sb.WriteString(p.Raw)

	return sb.String()
}

func AsOpenRouterError(resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	rd := bytes.NewReader(buf)

	var providerMeta errorResponse[ProviderError]

	err = json.NewDecoder(rd).Decode(&providerMeta)
	if err == nil {
		return &providerMeta.Error.Metadata
	}

	rd.Seek(0, io.SeekStart)

	var anyMeta errorResponse[any]

	err = json.NewDecoder(rd).Decode(&providerMeta)
	if err == nil {
		info := anyMeta.Error

		return &OpenRouterError{
			Message: info.Message,
			Code:    info.Code,
		}
	}

	return fmt.Errorf("http: %s", resp.Status)
}
