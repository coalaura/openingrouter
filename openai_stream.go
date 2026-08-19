package openingrouter

import "github.com/coalaura/openingrouter/internal/openai"

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
