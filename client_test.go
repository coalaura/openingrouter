package openingrouter

import (
	"context"
	"testing"
)

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
