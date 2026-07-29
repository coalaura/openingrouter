package openingrouter

import (
	"context"
)

// CreateSpeech synthesizes audio from the input text. The body of the returned
// response is not read or closed, closing it is up to the caller.
func (c *Client) CreateSpeech(ctx context.Context, request SpeechRequest) (*SpeechResponse, error) {
	req, err := c.NewRequest(ctx, "POST", "audio/speech", request)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return &SpeechResponse{
		ContentType: resp.Header.Get("Content-Type"),
		Body:        resp.Body,
	}, nil
}
