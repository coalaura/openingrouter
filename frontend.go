package openingrouter

import (
	"context"
	"encoding/json"
	"net/http"
)

func ListFrontendModels(ctx context.Context) ([]FrontendModel, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://openrouter.ai/api/frontend/v1/catalog/models", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result FrontendModelsResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}
