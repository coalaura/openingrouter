package openingrouter

import (
	"context"
	"encoding/json"
)

// ListImageModels retrieves the list of image generation models.
func (c *Client) ListImageModels(ctx context.Context) ([]ImageModel, error) {
	req, err := c.NewRequest(ctx, "GET", "images/models", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result OpenRouterResponse[[]ImageModel]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}
