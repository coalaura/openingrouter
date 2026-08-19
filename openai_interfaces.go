package openingrouter

import (
	"context"
	"net/http"
)

type OpenAICompatibleClient interface {
	NewRequest(ctx context.Context, method, path string, data any) (*http.Request, error)
	Do(req *http.Request) (*http.Response, error)

	ListModels(ctx context.Context, options *ListModelsOptions) ([]Model, error)

	CreateChatCompletion(ctx context.Context, request ChatCompletionRequest) (*ChatCompletionResponse, error)
	CreateChatCompletionStream(ctx context.Context, request ChatCompletionRequest) (OpenrouterStream[ChatStreamChunk], error)

	CreateCompletion(ctx context.Context, request CompletionRequest) (*CompletionResponse, error)
	CreateCompletionStream(ctx context.Context, request CompletionRequest) (OpenrouterStream[CompletionStreamChunk], error)

	CreateEmbeddings(ctx context.Context, request EmbeddingRequest) (*EmbeddingResponse, error)
}
