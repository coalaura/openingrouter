package openingrouter

import (
	"context"
	"encoding/json"
)

// CreateTranscription transcribes audio into text.
func (c *Client) CreateTranscription(ctx context.Context, request STTRequest) (*STTResponse, error) {
	req, err := c.NewRequest(ctx, "POST", "audio/transcriptions", request)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result STTResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
