package openingrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
