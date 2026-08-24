package openingrouter

import (
	"context"
	"encoding/json"
	"net/http"
)

// ListFrontendProviders retrieves the provider list from the OpenRouter frontend API.
func ListFrontendProviders(ctx context.Context) ([]FrontendProvider, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://openrouter.ai/api/frontend/v1/all-providers", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, AsOpenRouterError(resp, err)
	}

	defer resp.Body.Close()

	var result OpenRouterResponse[[]FrontendProvider]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}
