package openingrouter

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetModelBySlug retrieves detailed information about a model by its slug.
func (c *Client) GetModelBySlug(ctx context.Context, slug string) (*Model, error) {
	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("model/%s", slug), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result OpenRouterResponse[Model]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// ListModels retrieves the list of available models filtered by options.
func (c *Client) ListModels(ctx context.Context, options *ListModelsOptions) ([]Model, error) {
	req, err := c.NewRequest(ctx, "GET", "models", options)
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

// ListUserModels retrieves the list of models available to the current user.
func (c *Client) ListUserModels(ctx context.Context, options *ListUserModelsOptions) ([]Model, error) {
	req, err := c.NewRequest(ctx, "GET", "models/user", options)
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
