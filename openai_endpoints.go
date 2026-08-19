package openingrouter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coalaura/openingrouter/internal/openai"
)

// ListModels retrieves the list of models available on the OpenAI compatible API.
// Unlike the OpenRouter endpoint, the OpenAI list endpoint takes no query
// parameters, so no options are accepted.
func (c *OpenAIClient) ListModels(ctx context.Context) ([]Model, error) {
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

// CreateCompletion sends a legacy completions request and returns the response.
func (c *OpenAIClient) CreateCompletion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error) {
	if request.Stream == nil || *request.Stream {
		request.Stream = new(false)
	}

	req, err := c.NewRequest(ctx, "POST", "completions", completionRequestToOpenAI(&request))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result openai.CompletionResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return openaiCompletionResponseToResponse(&result), nil
}

// CreateCompletionStream sends a streaming completions request and returns a
// stream of completion chunks.
func (c *OpenAIClient) CreateCompletionStream(ctx context.Context, request CompletionRequest) (OpenrouterStream[CompletionStreamChunk], error) {
	if request.Stream == nil {
		request.Stream = new(true)
	}

	req, err := c.NewRequest(ctx, "POST", "completions", completionRequestToOpenAI(&request))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if IsResponseServerSentEventsStream(resp) {
		return &openAICompletionStream{
			inner: NewServerSentEventsStream[openai.CompletionChunk](ctx, resp),
		}, nil
	}

	defer resp.Body.Close()

	var fallback openai.CompletionResponse

	err = json.NewDecoder(resp.Body).Decode(&fallback)
	if err != nil {
		return nil, err
	}

	return NewJsonResponseStream(completionChunks(openaiCompletionResponseToResponse(&fallback))...), nil
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

// openAIChatStream adapts a stream of OpenAI wire chunks into a stream of
// openingrouter chat chunks, converting each element on receipt.
type openAIChatStream struct {
	inner OpenrouterStream[openai.ChatCompletionChunk]
}

// Recv returns the next converted chunk, or io.EOF when the stream is exhausted.
func (s *openAIChatStream) Recv() (ChatStreamChunk, error) {
	chunk, err := s.inner.Recv()
	if err != nil {
		var zero ChatStreamChunk

		return zero, err
	}

	return openaiChatCompletionChunkToChunk(&chunk), nil
}

// Close terminates the underlying stream.
func (s *openAIChatStream) Close() {
	s.inner.Close()
}

// openAICompletionStream adapts a stream of OpenAI wire chunks into a stream of
// openingrouter completion chunks, converting each element on receipt.
type openAICompletionStream struct {
	inner OpenrouterStream[openai.CompletionChunk]
}

// Recv returns the next converted chunk, or io.EOF when the stream is exhausted.
func (s *openAICompletionStream) Recv() (CompletionStreamChunk, error) {
	chunk, err := s.inner.Recv()
	if err != nil {
		var zero CompletionStreamChunk

		return zero, err
	}

	return openaiCompletionChunkToChunk(&chunk), nil
}

// Close terminates the underlying stream.
func (s *openAICompletionStream) Close() {
	s.inner.Close()
}

// completionChunks decomposes a non-streaming completions response into the
// chunks a stream would have yielded, for use as a fallback.
func completionChunks(response *CompletionResponse) []CompletionStreamChunk {
	if response == nil {
		return nil
	}

	chunk := CompletionStreamChunk{
		ID:                response.ID,
		Object:            CompletionObjectTextCompletion,
		Created:           response.Created,
		Model:             response.Model,
		SystemFingerprint: response.SystemFingerprint,
	}

	chunks := make([]CompletionStreamChunk, 0, len(response.Choices)+1)

	for _, choice := range response.Choices {
		finishReason := choice.FinishReason

		chunk.Choices = []CompletionStreamChoice{{
			Index:        choice.Index,
			Text:         choice.Text,
			FinishReason: &finishReason,
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
