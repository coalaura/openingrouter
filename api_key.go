package openingrouter

import (
	"context"
	"encoding/json"
)

// GetCurrentApiKey returns details about the API key used for the current session.
func (c *Client) GetCurrentApiKey(ctx context.Context) (*ApiKeyInfo, error) {
	req, err := c.NewRequest(ctx, "GET", "key", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result OpenRouterResponse[ApiKeyInfo]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}
