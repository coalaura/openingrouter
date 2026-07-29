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

type Client struct {
	client *http.Client
	base   *url.URL
	token  string

	title   string
	referer string
}

type Option func(*Client)

var defaultClient = &http.Client{}

func (c *Client) NewRequest(ctx context.Context, method, path string, data any) (*http.Request, error) {
	uri := c.base.ResolveReference(&url.URL{
		Path: path,
	})

	var (
		body        io.Reader
		contentType string
	)

	if data != nil {
		if method == "GET" {
			query, err := query.Values(data)
			if err != nil {
				return nil, err
			}

			uri.RawQuery = query.Encode()
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

	if c.referer != "" {
		req.Header.Set("HTTP-Referer", c.referer)
	}

	if c.title != "" {
		req.Header.Set("X-OpenRouter-Title", c.title)
	}

	return req, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, AsOpenRouterError(resp, err)
	}

	if resp.StatusCode > 299 {
		return nil, AsOpenRouterError(resp, nil)
	}

	return resp, nil
}

// NewClient returns a new client instance with the api token and options.
func NewClient(token string, options ...Option) *Client {
	client := &Client{
		client: defaultClient,
		base: &url.URL{
			Scheme: "https",
			Host:   "openrouter.ai",
			Path:   "/api/v1/",
		},
		token: fmt.Sprintf("Bearer %s", token),
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// WithClient sets the http client for each request.
func WithClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

// WithBase sets the base url of the client (default https://openrouter.ai/api/v1).
// Panics if the given base is not a valid url.
func WithBase(base string) Option {
	uri, err := url.Parse(base)
	if err == nil {
		panic(fmt.Sprintf("illegal base url: %v", err))
	}

	if !strings.HasSuffix(uri.Path, "/") {
		uri.Path += "/"
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
