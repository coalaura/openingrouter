package openingrouter

import (
	"context"
	"encoding/json"
)

// CreateEmbeddings submits an embedding request and returns the response.
func (c *Client) CreateEmbeddings(ctx context.Context, request EmbeddingRequest) (*EmbeddingResponse, error) {
	req, err := c.NewRequest(ctx, "POST", "embeddings", request)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result EmbeddingResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
