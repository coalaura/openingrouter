package openingrouter

import (
	"context"
	"encoding/json"
)

// GenerateImage generates images based on the provided request.
func (c *Client) GenerateImage(ctx context.Context, request ImageGenerationRequest) (*ImageGenerationResponse, error) {
	if request.Stream == nil || *request.Stream {
		request.Stream = new(false)
	}

	req, err := c.NewRequest(ctx, "POST", "images", request)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	var result ImageGenerationResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GenerateImageStream generates images as a stream of events.
func (c *Client) GenerateImageStream(ctx context.Context, request ImageGenerationRequest) (OpenrouterStream[ImageStreamEvent], error) {
	if request.Stream == nil {
		request.Stream = new(true)
	}

	req, err := c.NewRequest(ctx, "POST", "images", request)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if IsResponseServerSentEventsStream(resp) {
		return NewServerSentEventsStream[ImageStreamEvent](ctx, resp), nil
	}

	var fallback ImageGenerationResponse

	err = json.NewDecoder(resp.Body).Decode(&fallback)
	if err != nil {
		return nil, err
	}

	stream := NewJsonResponseStream[ImageStreamEvent]()

	for _, img := range fallback.Data {
		stream.Add(ImageStreamEvent{
			Type:      ImageStreamEventTypeCompleted,
			B64JSON:   img.B64JSON,
			MediaType: img.MediaType,
			Created:   fallback.Created,
			Usage:     fallback.Usage,
		})
	}

	return stream, nil
}
