package openingrouter

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) GetModelBySlug(ctx context.Context, slug string) (*Model, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("model/%s", slug), nil)
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

func (c *Client) ListModels(ctx context.Context, options *ListModelsOptions) ([]Model, error) {
	req, err := c.NewRequest("GET", "models", options)
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

func (c *Client) ListUserModels(ctx context.Context, options *ListUserModelsOptions) ([]Model, error) {
	req, err := c.NewRequest("GET", "models/user", options)
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
