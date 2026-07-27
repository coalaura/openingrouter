package openingrouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	base  *url.URL
	token string

	title   string
	referer string
}

type Option func(*Client)

func (c *Client) NewRequest(method, path string, data any) (*http.Request, error) {
	uri := c.base.ResolveReference(&url.URL{
		Path: path,
	})

	var body io.Reader

	if data != nil {
		var buf bytes.Buffer

		err := json.NewEncoder(&buf).Encode(data)
		if err != nil {
			return nil, err
		}

		body = &buf
	}

	req, err := http.NewRequest(method, uri.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.token)

	if c.referer != "" {
		req.Header.Set("HTTP-Referer", c.referer)
	}

	if c.title != "" {
		req.Header.Set("X-OpenRouter-Title", c.title)
	}

	return req, nil
}

// NewClient returns a new client instance with the api token and options.
func NewClient(token string, options ...Option) *Client {
	client := &Client{
		base: &url.URL{
			Scheme: "https",
			Host:   "openrouter.ai",
			Path:   "/api/v1",
		},
		token: fmt.Sprintf("Bearer %s", token),
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// WithBase sets the base url of the client (default https://openrouter.ai/api/v1)
// Panics if the given base is not a valid url.
func WithBase(base string) Option {
	uri, err := url.Parse(base)
	if err == nil {
		panic(fmt.Sprintf("illegal base url: %v", err))
	}

	uri.RawQuery = ""
	uri.Fragment = ""
	uri.RawFragment = ""

	return func(client *Client) {
		client.base = uri
	}
}

// WithTitle sets the X-Title header of each request.
func WithTitle(title string) Option {
	return func(client *Client) {
		client.title = title
	}
}

// WithReferer sets the HTTP-Referer header of each request.
func WithReferer(referer string) Option {
	return func(client *Client) {
		client.referer = referer
	}
}
