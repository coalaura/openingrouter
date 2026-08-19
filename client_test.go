package openingrouter

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient(
		"example-token",
		WithTitle("OpeningRouter"),
		WithReferer("https://example.com"),
		WithBase("https://example.com/api/v1"),
	)

	tAssertNotNil(t, client)
	tAssertEquals(t, client.title, "OpeningRouter")
	tAssertEquals(t, client.referer, "https://example.com")
	tAssertEquals(t, client.base.String(), "https://example.com/api/v1/")
}

func TestNewRequest(t *testing.T) {
	client := tCreateClient(t)

	req, err := client.NewRequest(context.Background(), "GET", "models", &ListModelsOptions{
		Limit: new(10),
	})

	tAssertNil(t, err)
	tAssertEquals(t, req.Method, "GET")
	tAssertEquals(t, req.URL.Path, "/api/v1/models")
	tAssertEquals(t, req.URL.RawQuery, "limit=10")
	tAssertEquals(t, req.URL.String(), "https://openrouter.ai/api/v1/models?limit=10")
}
