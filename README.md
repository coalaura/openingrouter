# openingrouter

Go client for the [OpenRouter](https://openrouter.ai) API.

## Install

```bash
go get github.com/coalaura/openingrouter
```

## Client

```go
client := openingrouter.NewClient(
    os.Getenv("OPENROUTER_API_KEY"),
    openingrouter.WithTitle("my-app"),
    openingrouter.WithReferer("https://example.com"),
)
```

Options: `WithClient`, `WithBase`, `WithTitle`, `WithReferer`.

## OpenAI-compatible endpoints

Point the same request/response types at any OpenAI-compatible API (OpenAI, Groq, Together, a local server, ...):

```go
oa := client.ToOpenAI() // reuses the token + HTTP client

// or standalone; defaults to https://api.openai.com/v1/
oa := openingrouter.NewOpenAIClient(
    os.Getenv("OPENAI_API_KEY"),
    openingrouter.WithOpenAIBase("https://api.groq.com/openai/v1"),
)

resp, err := oa.CreateChatCompletion(ctx, openingrouter.ChatCompletionRequest{...})
models, err := oa.ListModels(ctx)
```

Supported: `ListModels`, `GetModel`, `CreateChatCompletion`/`Stream`, `CreateEmbeddings`. Non-2xx responses map to `*OpenAIError`.

## Chat

```go
resp, err := client.CreateChatCompletion(ctx, openingrouter.ChatCompletionRequest{
    Model: "openai/gpt-oss-20b",
    Messages: []openingrouter.ChatMessage{{
        Role:    openingrouter.ChatRoleUser,
        Content: openingrouter.ChatContent{Text: "Hello"},
    }},
})
```

Streaming responses:

```go
stream, err := client.CreateChatCompletionStream(ctx, openingrouter.ChatCompletionRequest{
    Model: "openai/gpt-oss-20b",
    Messages: []openingrouter.ChatMessage{{
        Role:    openingrouter.ChatRoleUser,
        Content: openingrouter.ChatContent{Text: "Hello"},
    }},
})
if err != nil {
    return err
}

defer stream.Close()

for {
    chunk, err := stream.Recv()
    if errors.Is(err, io.EOF) {
        break
    }

    if err != nil {
        // *ChatStreamError when the stream fails after it started
        return err
    }

    // chunk.Choices[0].Delta.Content
}
```

Optional fields use pointers (`new(true)`, `new(0.7)`). Zero values and nil pointers are omitted from the request body. Set `MetadataLevel` to receive routing metadata.

## Embeddings

```go
resp, err := client.CreateEmbeddings(ctx, openingrouter.EmbeddingRequest{
    Model: "openai/text-embedding-3-small",
    Input: openingrouter.EmbeddingInput{Text: "Hello"},
})
```

`Input` accepts a single text, a list of texts, token arrays or multimodal content.

## Image generation

```go
resp, err := client.GenerateImage(ctx, openingrouter.ImageGenerationRequest{
    Model:       "black-forest-labs/flux.2-klein-4b",
    Prompt:      "a cat in a banana costume",
    AspectRatio: openingrouter.ImageAspectRatio1x1,
})

// resp.Data[i].B64JSON, resp.Data[i].MediaType
```

Streaming via `GenerateImageStream` (mid-stream failures return `*ImageStreamError` from `Recv`). List models with `ListImageModels`.

## Speech

Text-to-speech (caller owns and must close the body):

```go
resp, err := client.CreateSpeech(ctx, openingrouter.SpeechRequest{
    Model:          "sesame/csm-1b",
    Input:          "hello world",
    Voice:          "...",
    ResponseFormat: openingrouter.SpeechResponseFormatMP3,
})

defer resp.Body.Close()

// resp.Body, resp.ContentType, resp.GenerationID
```

Speech-to-text:

```go
resp, err := client.CreateTranscription(ctx, openingrouter.STTRequest{
    Model: "google/chirp-3",
    InputAudio: openingrouter.STTInputAudio{
        Data:   base64Audio,
        Format: "wav",
    },
})

// resp.Text
```

## Models

```go
models, err := client.ListModels(ctx, &openingrouter.ListModelsOptions{
    Limit: new(10),
})

model, err := client.GetModelBySlug(ctx, "deepseek/deepseek-v4-flash")

userModels, err := client.ListUserModels(ctx, nil)
embeddingModels, err := client.ListEmbeddingModels(ctx, nil)
```

## API key

```go
info, err := client.GetCurrentApiKey(ctx)
```

## Frontend catalog

Unauthenticated frontend route (not part of the public API, sometimes has more information):

```go
models, err := openingrouter.ListFrontendModels(ctx)
```

## Errors

`Client.Do` maps non-2xx responses to:

| Type | When |
|------|------|
| `*OpenRouterError` | API error with numeric code |
| `*ApiError` | Named API / validation error |
| `*ProviderError` | Upstream provider error |

Network failures are returned as-is.

Streaming endpoints that fail after the response has started surface the failure from `Recv` as a typed error (not as a field on the chunk):

| Type | Endpoint |
|------|----------|
| `*ChatStreamError` | `CreateChatCompletionStream` |
| `*ImageStreamError` | `GenerateImageStream` |

Their `Error()` strings match the same style as `*OpenRouterError` (`openrouter code <code>: <message>`, falling back to `openrouter: <message>` when no code is present). Nested provider JSON in the message is cleaned the same way as HTTP errors. Inspect fields with `errors.As` / `errors.AsType`:

```go
chunk, err := stream.Recv()
if err != nil {
    if errors.Is(err, io.EOF) {
        break
    }

    if streamErr, ok := errors.AsType[*openingrouter.ChatStreamError](err); ok {
        // streamErr.Code, streamErr.Metadata
        return streamErr
    }

    return err
}
```

## Streams

`OpenrouterStream[T]` is the common interface:

```go
type OpenrouterStream[T any] interface {
    Recv() (T, error) // io.EOF when done; typed stream error on mid-stream failure
    Close()
}
```

Always `defer stream.Close()`. Chunk types that carry an in-band error implement it privately; `Recv` promotes that error so callers only need the usual `err != nil` path.

## Layout

| File prefix | Endpoint |
|-------------|----------|
| `chat_` | `POST /chat/completions` |
| `embeddings_` | `POST /embeddings` |
| `embedding_models_` | `GET /embeddings/models` |
| `image_` | `POST /images` |
| `image_models_` | `GET /images/models` |
| `models_` | `GET /models`, `/models/user`, `/model/{slug}` |
| `stt_` | `POST /audio/transcriptions` |
| `tts_` | `POST /audio/speech` |
| `api_key_` | `GET /key` |
| `frontend_` | frontend catalog |
| `common_types.go` | shared request/response types |

## Tests

Integration tests hit the live API. Set `OPENROUTER_API_KEY` or they skip:

```bash
OPENROUTER_API_KEY="sk-or-v1-46...2f" go test ./...
```