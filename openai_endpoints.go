package openingrouter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coalaura/openingrouter/internal/openai"
)

// ListModels retrieves the list of models available on the OpenAI compatible API.
func (c *OpenAIClient) ListModels(ctx context.Context, _ *ListModelsOptions) ([]Model, error) {
	req, err := c.NewRequest(ctx, "GET", "models", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result openai.ModelsList

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return openaiModelsListToModels(&result), nil
}

// GetModel retrieves a single model by its id.
func (c *OpenAIClient) GetModel(ctx context.Context, model string) (*Model, error) {
	req, err := c.NewRequest(ctx, "GET", fmt.Sprintf("models/%s", model), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result openai.Model

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return openaiModelToModel(&result), nil
}

// CreateChatCompletion sends a chat completion request and returns the response.
func (c *OpenAIClient) CreateChatCompletion(ctx context.Context, request ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if request.Stream == nil || *request.Stream {
		request.Stream = new(false)
	}

	req, err := c.NewRequest(ctx, "POST", "chat/completions", chatCompletionRequestToOpenAI(&request))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result openai.ChatCompletionResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return openaiChatCompletionResponseToResponse(&result), nil
}

// CreateChatCompletionStream sends a streaming chat completion request and
// returns a stream of completion chunks.
func (c *OpenAIClient) CreateChatCompletionStream(ctx context.Context, request ChatCompletionRequest) (OpenrouterStream[ChatStreamChunk], error) {
	if request.Stream == nil {
		request.Stream = new(true)
	}

	req, err := c.NewRequest(ctx, "POST", "chat/completions", chatCompletionRequestToOpenAI(&request))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if IsResponseServerSentEventsStream(resp) {
		return &openAIChatStream{
			inner: NewServerSentEventsStream[openai.ChatCompletionChunk](ctx, resp),
		}, nil
	}

	defer resp.Body.Close()

	var fallback openai.ChatCompletionResponse

	err = json.NewDecoder(resp.Body).Decode(&fallback)
	if err != nil {
		return nil, err
	}

	return NewJsonResponseStream(chatCompletionChunks(openaiChatCompletionResponseToResponse(&fallback))...), nil
}

// CreateEmbeddings submits an embedding request and returns the response.
func (c *OpenAIClient) CreateEmbeddings(ctx context.Context, request EmbeddingRequest) (*EmbeddingResponse, error) {
	req, err := c.NewRequest(ctx, "POST", "embeddings", embeddingRequestToOpenAI(&request))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result openai.EmbeddingResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return openaiEmbeddingResponseToResponse(&result), nil
}

