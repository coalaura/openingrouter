package openai

// CompletionObject is the object type of a completion response or chunk.
type CompletionObject string

const (
	CompletionObjectTextCompletion CompletionObject = "text_completion"
)

// CompletionInput represents text, token, or batch prompt input(s). Exactly one
// of the fields should be set. TokenArrays takes precedence over Tokens, Texts
// and Text when more than one is set.
type CompletionInput struct {
	Text        string
	Texts       []string
	Tokens      []int
	TokenArrays [][]int
}

// MarshalJSON implements the json.Marshaler interface for CompletionInput.
func (ci CompletionInput) MarshalJSON() ([]byte, error) {
	switch {
	case len(ci.TokenArrays) > 0:
		return jsonMarshal(ci.TokenArrays)
	case len(ci.Tokens) > 0:
		return jsonMarshal(ci.Tokens)
	case len(ci.Texts) > 0:
		return jsonMarshal(ci.Texts)
	default:
		return jsonMarshal(ci.Text)
	}
}

// CompletionRequest is the request body of the completions endpoint. Model and
// Prompt are required, zero values and nil pointers of the remaining fields are
// omitted.
type CompletionRequest struct {
	Model  string          `json:"model"`
	Prompt CompletionInput `json:"prompt"`

	BestOf           *int             `json:"best_of,omitempty"`
	Echo             *bool            `json:"echo,omitempty"`
	FrequencyPenalty *float64         `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
	Logprobs         *int             `json:"logprobs,omitempty"`
	MaxTokens        *int             `json:"max_tokens,omitempty"`
	N                *int             `json:"n,omitempty"`
	PresencePenalty  *float64         `json:"presence_penalty,omitempty"`
	Seed             *int             `json:"seed,omitempty"`
	Stop             []string         `json:"stop,omitempty"`
	Stream           *bool            `json:"stream,omitempty"`
	StreamOptions    *ChatStreamOptions `json:"stream_options,omitempty"`
	Suffix           string           `json:"suffix,omitempty"`
	Temperature      *float64         `json:"temperature,omitempty"`
	TopP             *float64         `json:"top_p,omitempty"`
	User             string           `json:"user,omitempty"`
}

// CompletionResponse is the root response of the completions endpoint.
type CompletionResponse struct {
	ID                string             `json:"id"`
	Object            CompletionObject   `json:"object"`
	Created           int64              `json:"created"`
	Model             string             `json:"model"`
	Choices           []CompletionChoice `json:"choices"`
	SystemFingerprint *string            `json:"system_fingerprint"`
	Usage             *CompletionUsage   `json:"usage,omitempty"`
}

// CompletionChoice represents a single completion choice.
type CompletionChoice struct {
	Index        int              `json:"index"`
	FinishReason string           `json:"finish_reason"`
	Text         string           `json:"text"`
	Logprobs     *CompletionLogprobs `json:"logprobs,omitempty"`
}

// CompletionLogprobs represents the log probabilities of a completion choice.
type CompletionLogprobs struct {
	TextOffset    []int              `json:"text_offset,omitempty"`
	TokenLogprobs []float64          `json:"token_logprobs,omitempty"`
	Tokens        []string           `json:"tokens,omitempty"`
	TopLogprobs   []map[string]float64 `json:"top_logprobs,omitempty"`
}

// CompletionUsage represents the token usage of a completion.
type CompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// CompletionChunk represents a single chunk of a streaming completion.
type CompletionChunk struct {
	ID                string                `json:"id"`
	Object            CompletionObject      `json:"object"`
	Created           int64                 `json:"created"`
	Model             string                `json:"model"`
	Choices           []CompletionChunkChoice `json:"choices"`
	SystemFingerprint *string               `json:"system_fingerprint"`
	Usage             *CompletionUsage      `json:"usage,omitempty"`
}

// CompletionChunkChoice represents a single choice of a streaming completion.
type CompletionChunkChoice struct {
	Index        int              `json:"index"`
	Text         string           `json:"text"`
	FinishReason *string          `json:"finish_reason"`
	Logprobs     *CompletionLogprobs `json:"logprobs,omitempty"`
}
