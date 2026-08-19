package openingrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
)

// OpenAIClient represents an OpenAI compatible API client. It sends requests in
// the OpenAI wire format and exposes the same methods as Client, returning the
// openingrouter types. Use NewOpenAIClient to build one directly or Client.ToOpenAI
// to derive one from an existing OpenRouter client.
type OpenAIClient struct {
	client *http.Client
	base   *url.URL
	token  string
}

// OpenAIOption configures an OpenAI compatible API client.
type OpenAIOption func(*OpenAIClient)

// NewRequest constructs a new HTTP request targeting the OpenAI compatible API.
// It mirrors Client.NewRequest but omits the OpenRouter specific headers.
func (c *OpenAIClient) NewRequest(ctx context.Context, method, path string, data any) (*http.Request, error) {
	uri := c.base.ResolveReference(&url.URL{
		Path: path,
	})

	var (
		body        io.Reader
		contentType string
	)

	if data != nil {
		if method == "GET" {
			values, err := query.Values(data)
			if err != nil {
				return nil, err
			}

			uri.RawQuery = values.Encode()
		} else {
			var buf bytes.Buffer

			err := json.NewEncoder(&buf).Encode(data)
			if err != nil {
				return nil, err
			}

			body = &buf

			contentType = "application/json"
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, uri.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.token)
	req.Header.Set("User-Agent", "openingrouter")

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return req, nil
}

// Do sends an HTTP request and processes any returned API errors.
func (c *OpenAIClient) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, AsOpenAIError(resp, err)
	}

	if resp.StatusCode > 299 {
		return nil, AsOpenAIError(resp, nil)
	}

	return resp, nil
}

// NewOpenAIClient returns a new OpenAI compatible client instance with the api
// token and options. The default base url is https://api.openai.com/v1.
func NewOpenAIClient(token string, options ...OpenAIOption) *OpenAIClient {
	client := &OpenAIClient{
		client: defaultClient,
		base: &url.URL{
			Scheme: "https",
			Host:   "api.openai.com",
			Path:   "/v1/",
		},
		token: fmt.Sprintf("Bearer %s", token),
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// ToOpenAI returns an OpenAI compatible client derived from this OpenRouter
// client, carrying over the http client and api token. The given options are
// applied on top, so the base url and http client can be overridden.
func (c *Client) ToOpenAI(options ...OpenAIOption) *OpenAIClient {
	token := strings.TrimPrefix(c.token, "Bearer ")

	client := NewOpenAIClient(token, options...)

	if c.client != nil {
		client.client = c.client
	}

	return client
}

// WithOpenAIBase sets the base url of the client (default https://api.openai.com/v1).
// Panics if the given base is not a valid url.
func WithOpenAIBase(base string) OpenAIOption {
	uri, err := url.Parse(base)
	if err != nil {
		panic(fmt.Sprintf("illegal base url: %v", err))
	}

	if !strings.HasSuffix(uri.Path, "/") {
		uri.Path += "/"
	}

	uri.RawQuery = ""
	uri.Fragment = ""
	uri.RawFragment = ""

	return func(c *OpenAIClient) {
		c.base = uri
	}
}

// WithOpenAIHTTPClient sets the http client for each request.
func WithOpenAIHTTPClient(client *http.Client) OpenAIOption {
	return func(c *OpenAIClient) {
		c.client = client
	}
}

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

	if decodeErr := json.NewDecoder(resp.Body).Decode(&body); decodeErr != nil ||
		body.Error.Message == "" && body.Error.Type == "" && len(body.Error.Code) == 0 {
		return newHttpError(resp)
	}

	info := body.Error

	var code int64

	if raw := strings.TrimSpace(string(info.Code)); raw != "" && raw != "null" {
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
