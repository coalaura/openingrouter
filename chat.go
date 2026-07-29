package openingrouter

import (
	"context"
	"encoding/json"
)

func (c *Client) CreateChatCompletion(ctx context.Context, request ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if request.Stream == nil || *request.Stream {
		request.Stream = new(false)
	}

	req, err := c.NewRequest("POST", "chat/completions", request)
	if err != nil {
		return nil, err
	}

	if request.MetadataLevel != "" {
		req.Header.Set("X-OpenRouter-Metadata", string(request.MetadataLevel))
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result ChatCompletionResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) CreateChatCompletionStream(ctx context.Context, request ChatCompletionRequest) (OpenrouterStream[ChatStreamChunk], error) {
	if request.Stream == nil {
		request.Stream = new(true)
	}

	req, err := c.NewRequest("POST", "chat/completions", request)
	if err != nil {
		return nil, err
	}

	if request.MetadataLevel != "" {
		req.Header.Set("X-OpenRouter-Metadata", string(request.MetadataLevel))
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if IsResponseServerSentEventsStream(resp) {
		return NewServerSentEventsStream[ChatStreamChunk](ctx, resp), nil
	}

	defer resp.Body.Close()

	var fallback ChatCompletionResponse

	err = json.NewDecoder(resp.Body).Decode(&fallback)
	if err != nil {
		return nil, err
	}

	return NewJsonResponseStream(chatCompletionChunks(&fallback)...), nil
}

func chatCompletionChunks(response *ChatCompletionResponse) []ChatStreamChunk {
	chunk := ChatStreamChunk{
		ID:                 response.ID,
		Object:             ChatObjectCompletionChunk,
		Created:            response.Created,
		Model:              response.Model,
		ServiceTier:        response.ServiceTier,
		OpenRouterMetadata: response.OpenRouterMetadata,
	}

	if response.SystemFingerprint != nil {
		chunk.SystemFingerprint = *response.SystemFingerprint
	}

	// metadata chunk, optional usage chunk and roughly four deltas per choice
	chunks := make([]ChatStreamChunk, 0, 2+len(response.Choices)*4)

	for _, choice := range response.Choices {
		chunk.Choices = append(chunk.Choices, ChatStreamChoice{
			Index: choice.Index,
			Delta: ChatStreamDelta{
				Role: choice.Message.Role,
			},
		})
	}

	chunks = append(chunks, chunk)

	chunk.OpenRouterMetadata = nil
	chunk.Choices = nil

	emit := func(index int, delta ChatStreamDelta) {
		chunk.Choices = []ChatStreamChoice{{
			Index: index,
			Delta: delta,
		}}

		chunks = append(chunks, chunk)
	}

	for _, choice := range response.Choices {
		for _, detail := range choice.Message.ReasoningDetails {
			emit(choice.Index, ChatStreamDelta{
				ReasoningDetails: []ChatReasoningDetail{detail},
			})
		}

		if choice.Message.Reasoning != "" {
			emit(choice.Index, ChatStreamDelta{
				Reasoning: choice.Message.Reasoning,
			})
		}

		content := choice.Message.Content.String()
		if content != "" {
			emit(choice.Index, ChatStreamDelta{
				Content: content,
			})
		}

		if choice.Message.Refusal != "" {
			emit(choice.Index, ChatStreamDelta{
				Refusal: choice.Message.Refusal,
			})
		}

		if choice.Message.Audio != nil {
			emit(choice.Index, ChatStreamDelta{
				Audio: choice.Message.Audio,
			})
		}

		for _, image := range choice.Message.Images {
			emit(choice.Index, ChatStreamDelta{
				Images: []ChatAssistantImage{image},
			})
		}

		for index, call := range choice.Message.ToolCalls {
			emit(choice.Index, ChatStreamDelta{
				ToolCalls: []ChatStreamToolCall{{
					Index: index,
					ID:    call.ID,
					Type:  call.Type,
					Function: &ChatStreamToolCallFunction{
						Name:      call.Function.Name,
						Arguments: call.Function.Arguments,
					},
				}},
			})
		}

		chunk.Choices = []ChatStreamChoice{{
			Index:        choice.Index,
			FinishReason: choice.FinishReason,
			Logprobs:     choice.Logprobs,
		}}

		chunks = append(chunks, chunk)
	}

	if response.Usage != nil {
		chunk.Choices = nil
		chunk.Usage = response.Usage

		chunks = append(chunks, chunk)
	}

	return chunks
}
