package openingrouter

import (
	"context"
	"encoding/json"
)

// ListEmbeddingModels retrieves the list of available embeddings models.
func (c *Client) ListEmbeddingModels(ctx context.Context, options *ListEmbeddingModelsOptions) ([]Model, error) {
	req, err := c.NewRequest(ctx, "GET", "embeddings/models", options)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result OpenRouterResponse[[]Model]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}
