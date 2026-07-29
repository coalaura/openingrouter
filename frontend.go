package openingrouter

import (
	"context"
	"encoding/json"
	"net/http"
)

// ListFrontendModels retrieves the model catalog from the OpenRouter frontend API.
func ListFrontendModels(ctx context.Context) ([]FrontendModel, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://openrouter.ai/api/frontend/v1/catalog/models", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, AsOpenRouterError(resp, err)
	}

	defer resp.Body.Close()

	var result OpenRouterResponse[[]FrontendModel]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}
