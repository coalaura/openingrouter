package openingrouter

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// tCreateOpenAIClient builds an OpenAIClient backed by a local test server.
func tCreateOpenAIClient(t testing.TB, handler http.HandlerFunc) *OpenAIClient {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return NewOpenAIClient("test-token", WithOpenAIBase(server.URL))
}

// tAssertJSONString asserts body[key] is a string equal to expected.
func tAssertJSONString(t testing.TB, body map[string]any, key, expected string) {
	t.Helper()

	value, ok := body[key].(string)
	if !ok || value != expected {
		t.Fatalf("expected body[%s] = %q, got %v", key, expected, body[key])
	}
}

// tAssertJSONBool asserts body[key] is a bool equal to expected.
func tAssertJSONBool(t testing.TB, body map[string]any, key string, expected bool) {
	t.Helper()

	value, ok := body[key].(bool)
	if !ok || value != expected {
		t.Fatalf("expected body[%s] = %v, got %v", key, expected, body[key])
	}
}

// tDecodeBody decodes the request body into a generic json object.
func tDecodeBody(t testing.TB, r *http.Request) map[string]any {
	t.Helper()

	var body map[string]any

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Fatalf("decode request body: %v", err)
	}

	return body
}

func TestOpenAIListModels(t *testing.T) {
	client := tCreateOpenAIClient(t, func(w http.ResponseWriter, r *http.Request) {
		tAssertEquals(t, r.Method, "GET")
		tAssertEquals(t, r.URL.Path, "/models")
		tAssertEquals(t, r.Header.Get("Authorization"), "Bearer test-token")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"object": "list",
			"data": [
				{"id": "gpt-4o", "object": "model", "created": 1687882411, "owned_by": "openai"},
				{"id": "gpt-3.5-turbo", "object": "model", "created": 1687882411, "owned_by": "openai"}
			]
		}`))
	})

	models, err := client.ListModels(context.Background())

	tAssertNil(t, err)
	tAssertLen(t, models, 2)
	tAssertEquals(t, models[0].ID, "gpt-4o")
	tAssertEquals(t, models[0].CanonicalSlug, "gpt-4o")
	tAssertEquals(t, models[0].Name, "gpt-4o")
	tAssertEquals(t, models[0].Created, int64(1687882411))
}

func TestOpenAIGetModel(t *testing.T) {
	client := tCreateOpenAIClient(t, func(w http.ResponseWriter, r *http.Request) {
		tAssertEquals(t, r.Method, "GET")
		tAssertEquals(t, r.URL.Path, "/models/gpt-4o")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"id": "gpt-4o",
			"object": "model",
			"created": 1687882411,
			"owned_by": "openai",
			"shutdown_date": "2025-12-31"
		}`))
	})

	model, err := client.GetModel(context.Background(), "gpt-4o")

	tAssertNil(t, err)
	tAssertEquals(t, model.ID, "gpt-4o")
	tAssertEquals(t, model.Name, "gpt-4o")
	tAssertNotNil(t, model.ExpirationDate)
}

func TestOpenAICreateChatCompletion(t *testing.T) {
	client := tCreateOpenAIClient(t, func(w http.ResponseWriter, r *http.Request) {
		tAssertEquals(t, r.Method, "POST")
		tAssertEquals(t, r.URL.Path, "/chat/completions")

		body := tDecodeBody(t, r)
		tAssertJSONString(t, body, "model", "gpt-4o")
		tAssertJSONBool(t, body, "stream", false)

		messages, ok := body["messages"].([]any)
		if !ok || len(messages) != 1 {
			t.Fatalf("expected 1 message, got %v", body["messages"])
		}

		first, ok := messages[0].(map[string]any)
		if !ok {
			t.Fatalf("expected message object, got %T", messages[0])
		}

		tAssertJSONString(t, first, "role", "user")
		tAssertJSONString(t, first, "content", "hi")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"id": "chatcmpl-123",
			"object": "chat.completion",
			"created": 1677652288,
			"model": "gpt-4o",
			"choices": [{
				"index": 0,
				"finish_reason": "stop",
				"message": {"role": "assistant", "content": "hello world"}
			}],
			"usage": {"prompt_tokens": 9, "completion_tokens": 12, "total_tokens": 21}
		}`))
	})

	resp, err := client.CreateChatCompletion(context.Background(), ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: ChatRoleUser, Content: ChatContent{Text: "hi"}},
		},
	})

	tAssertNil(t, err)
	tAssertEquals(t, resp.ID, "chatcmpl-123")
	tAssertEquals(t, resp.Object, ChatObjectCompletion)
	tAssertLen(t, resp.Choices, 1)
	tAssertEquals(t, resp.Choices[0].Message.Role, ChatRoleAssistant)
	tAssertEquals(t, resp.Choices[0].Message.Content.Text, "hello world")
	tAssertNotNil(t, resp.Usage)
	tAssertEquals(t, resp.Usage.TotalTokens, 21)
}

func TestOpenAICreateChatCompletionStream(t *testing.T) {
	client := tCreateOpenAIClient(t, func(w http.ResponseWriter, r *http.Request) {
		tAssertEquals(t, r.URL.Path, "/chat/completions")

		body := tDecodeBody(t, r)
		tAssertJSONBool(t, body, "stream", true)

		w.Header().Set("Content-Type", "text/event-stream")

		flusher, _ := w.(http.Flusher)

		w.Write([]byte(`data: {"id":"chatcmpl-1","object":"chat.completion.chunk","created":1,"model":"gpt-4o","choices":[{"index":0,"delta":{"role":"assistant","content":"hello"}}]}` + "\n\n"))
		if flusher != nil {
			flusher.Flush()
		}

		w.Write([]byte(`data: {"id":"chatcmpl-1","object":"chat.completion.chunk","created":1,"model":"gpt-4o","choices":[{"index":0,"delta":{"content":" world"},"finish_reason":"stop"}]}` + "\n\n"))
		if flusher != nil {
			flusher.Flush()
		}

		w.Write([]byte("data: [DONE]" + "\n\n"))
		if flusher != nil {
			flusher.Flush()
		}
	})

	stream, err := client.CreateChatCompletionStream(context.Background(), ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: ChatRoleUser, Content: ChatContent{Text: "hi"}},
		},
	})

	tAssertNil(t, err)
	defer stream.Close()

	var (
		role   ChatRole
		result strings.Builder
	)

	for {
		chunk, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		tAssertNil(t, err)

		if len(chunk.Choices) == 0 {
			continue
		}

		delta := chunk.Choices[0].Delta
		if role == "" && delta.Role != "" {
			role = delta.Role
		}

		result.WriteString(delta.Content)
	}

	tAssertEquals(t, role, ChatRoleAssistant)
	tAssertEquals(t, result.String(), "hello world")
}

func TestOpenAICreateCompletion(t *testing.T) {
	client := tCreateOpenAIClient(t, func(w http.ResponseWriter, r *http.Request) {
		tAssertEquals(t, r.Method, "POST")
		tAssertEquals(t, r.URL.Path, "/completions")

		body := tDecodeBody(t, r)
		tAssertJSONString(t, body, "model", "gpt-3.5-turbo")
		tAssertJSONString(t, body, "prompt", "Say riff raff")
		tAssertJSONBool(t, body, "stream", false)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"id": "cmpl-123",
			"object": "text_completion",
			"created": 1677652288,
			"model": "gpt-3.5-turbo",
			"choices": [{"index": 0, "text": "riff raff", "finish_reason": "stop"}],
			"usage": {"prompt_tokens": 5, "completion_tokens": 3, "total_tokens": 8}
		}`))
	})

	resp, err := client.CreateCompletion(context.Background(), CompletionRequest{
		Model:  "gpt-3.5-turbo",
		Prompt: CompletionInput{Text: "Say riff raff"},
	})

	tAssertNil(t, err)
	tAssertEquals(t, resp.ID, "cmpl-123")
	tAssertEquals(t, resp.Object, CompletionObjectTextCompletion)
	tAssertLen(t, resp.Choices, 1)
	tAssertEquals(t, resp.Choices[0].Text, "riff raff")
	tAssertNotNil(t, resp.Usage)
	tAssertEquals(t, resp.Usage.TotalTokens, 8)
}

func TestOpenAICreateCompletionStream(t *testing.T) {
	client := tCreateOpenAIClient(t, func(w http.ResponseWriter, r *http.Request) {
		tAssertEquals(t, r.URL.Path, "/completions")

		body := tDecodeBody(t, r)
		tAssertJSONBool(t, body, "stream", true)

		w.Header().Set("Content-Type", "text/event-stream")

		w.Write([]byte(`data: {"id":"cmpl-1","object":"text_completion","created":1,"model":"gpt-3.5-turbo","choices":[{"index":0,"text":"riff"}]}` + "\n\n"))
		w.Write([]byte(`data: {"id":"cmpl-1","object":"text_completion","created":1,"model":"gpt-3.5-turbo","choices":[{"index":0,"text":" raff","finish_reason":"stop"}]}` + "\n\n"))
		w.Write([]byte("data: [DONE]" + "\n\n"))
	})

	stream, err := client.CreateCompletionStream(context.Background(), CompletionRequest{
		Model:  "gpt-3.5-turbo",
		Prompt: CompletionInput{Text: "Say riff raff"},
	})

	tAssertNil(t, err)
	defer stream.Close()

	var result strings.Builder

	for {
		chunk, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		tAssertNil(t, err)

		if len(chunk.Choices) == 0 {
			continue
		}

		result.WriteString(chunk.Choices[0].Text)
	}

	tAssertEquals(t, result.String(), "riff raff")
}

func TestOpenAICreateEmbeddings(t *testing.T) {
	client := tCreateOpenAIClient(t, func(w http.ResponseWriter, r *http.Request) {
		tAssertEquals(t, r.Method, "POST")
		tAssertEquals(t, r.URL.Path, "/embeddings")

		body := tDecodeBody(t, r)
		tAssertJSONString(t, body, "model", "text-embedding-3-small")
		tAssertJSONString(t, body, "input", "hello")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"object": "list",
			"data": [{"object": "embedding", "embedding": [0.1, 0.2, 0.3], "index": 0}],
			"model": "text-embedding-3-small",
			"usage": {"prompt_tokens": 8, "total_tokens": 8}
		}`))
	})

	resp, err := client.CreateEmbeddings(context.Background(), EmbeddingRequest{
		Model: "text-embedding-3-small",
		Input: EmbeddingInput{Text: "hello"},
	})

	tAssertNil(t, err)
	tAssertEquals(t, resp.Object, EmbeddingObjectList)
	tAssertLen(t, resp.Data, 1)
	tAssertLen(t, resp.Data[0].Embedding.Floats, 3)
	tAssertEquals(t, resp.Data[0].Embedding.Floats[0], 0.1)
	tAssertEquals(t, resp.Data[0].Index, 0)
	tAssertNotNil(t, resp.Usage)
	tAssertEquals(t, resp.Usage.TotalTokens, 8)
}

func TestOpenAIError(t *testing.T) {
	client := tCreateOpenAIClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"message":"Invalid model","type":"invalid_request_error","code":"model_not_found"}}`))
	})

	_, err := client.CreateChatCompletion(context.Background(), ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: ChatRoleUser, Content: ChatContent{Text: "hi"}},
		},
	})

	tAssertNotNil(t, err)

	var openAIErr *OpenAIError
	if !errors.As(err, &openAIErr) {
		t.Fatalf("expected *OpenAIError, got %T: %v", err, err)
	}

	tAssertEquals(t, openAIErr.Message, "Invalid model")
	tAssertEquals(t, openAIErr.Type, "invalid_request_error")
	tAssertEquals(t, openAIErr.Code, "model_not_found")
	tAssertEquals(t, openAIErr.ErrorStatus.Code, int64(400))
}
