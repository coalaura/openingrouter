package openingrouter

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// EmbeddingRequest represents the request body of the embeddings endpoint.
// Model and Input are required, zero values and nil pointers of the remaining
// fields are omitted from the request.
type EmbeddingRequest struct {
	Model string         `json:"model"`
	Input EmbeddingInput `json:"input"`

	Dimensions     *int                    `json:"dimensions,omitempty"`
	EncodingFormat EmbeddingEncodingFormat `json:"encoding_format,omitempty"`
	InputType      string                  `json:"input_type,omitempty"`
	Provider       *ProviderPreferences    `json:"provider,omitempty"`
	User           string                  `json:"user,omitempty"`
}

// EmbeddingInput represents text, token, or multimodal input(s) to embed. Exactly
// one of the fields should be set. MultimodalInputs takes precedence over
// TokenArrays, Tokens, Texts and Text when more than one is set.
type EmbeddingInput struct {
	Text             string
	Texts            []string
	Tokens           []int
	TokenArrays      [][]int
	MultimodalInputs []EmbeddingMultimodalInput
}

// MarshalJSON implements the json.Marshaler interface for EmbeddingInput.
func (ei EmbeddingInput) MarshalJSON() ([]byte, error) {
	switch {
	case len(ei.MultimodalInputs) > 0:
		return json.Marshal(ei.MultimodalInputs)
	case len(ei.TokenArrays) > 0:
		return json.Marshal(ei.TokenArrays)
	case len(ei.Tokens) > 0:
		return json.Marshal(ei.Tokens)
	case len(ei.Texts) > 0:
		return json.Marshal(ei.Texts)
	default:
		return json.Marshal(ei.Text)
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface for EmbeddingInput.
func (ei *EmbeddingInput) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		*ei = EmbeddingInput{}

		return nil
	}

	if len(data) > 0 && data[0] == '"' {
		return json.Unmarshal(data, &ei.Text)
	}

	if len(data) == 0 || data[0] != '[' {
		return fmt.Errorf("embedding input: unexpected json value")
	}

	trimmed := bytes.TrimSpace(data[1:])
	if len(trimmed) == 0 || trimmed[0] == ']' {
		ei.Texts = []string{}

		return nil
	}

	switch trimmed[0] {
	case '"':
		return json.Unmarshal(data, &ei.Texts)
	case '[':
		return json.Unmarshal(data, &ei.TokenArrays)
	case '{':
		return json.Unmarshal(data, &ei.MultimodalInputs)
	default:
		return json.Unmarshal(data, &ei.Tokens)
	}
}

// EmbeddingMultimodalInput represents a single multimodal embedding input made
// of one or more content parts.
type EmbeddingMultimodalInput struct {
	Content []EmbeddingContentPart `json:"content"`
}

// EmbeddingContentPart represents a single content part of a multimodal
// embedding input. Type determines which of the remaining fields are used.
type EmbeddingContentPart struct {
	Type       EmbeddingContentPartType `json:"type"`
	Text       string                   `json:"text,omitempty"`
	ImageURL   *ContentPartImageURL     `json:"image_url,omitempty"`
	InputAudio *EmbeddingMedia          `json:"input_audio,omitempty"`
	InputVideo *EmbeddingMedia          `json:"input_video,omitempty"`
	InputFile  *EmbeddingMedia          `json:"input_file,omitempty"`
}

// EmbeddingMedia holds base64-encoded media and its format for an audio, video
// or file embedding content part.
type EmbeddingMedia struct {
	Data   string `json:"data"`
	Format string `json:"format,omitempty"`
}

// EmbeddingResponse is the root response of the embeddings endpoint.
type EmbeddingResponse struct {
	ID     string          `json:"id,omitempty"`
	Object EmbeddingObject `json:"object"`
	Data   []Embedding     `json:"data"`
	Model  string          `json:"model"`
	Usage  *EmbeddingUsage `json:"usage,omitempty"`
}

// Embedding represents a single embedding object of an embeddings response.
type Embedding struct {
	Object    EmbeddingObject `json:"object"`
	Embedding EmbeddingValue  `json:"embedding"`
	Index     int             `json:"index"`
}

// EmbeddingValue represents an embedding vector encoded either as an array of
// floats or as a base64 string. Floats takes precedence over Base64 when both
// are set.
type EmbeddingValue struct {
	Floats []float64
	Base64 string
}

// MarshalJSON implements the json.Marshaler interface for EmbeddingValue.
func (ev EmbeddingValue) MarshalJSON() ([]byte, error) {
	if len(ev.Floats) > 0 || ev.Base64 == "" {
		if ev.Floats == nil {
			return []byte("[]"), nil
		}

		return json.Marshal(ev.Floats)
	}

	return json.Marshal(ev.Base64)
}

// UnmarshalJSON implements the json.Unmarshaler interface for EmbeddingValue.
func (ev *EmbeddingValue) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		*ev = EmbeddingValue{}

		return nil
	}

	if len(data) > 0 && data[0] == '"' {
		return json.Unmarshal(data, &ev.Base64)
	}

	return json.Unmarshal(data, &ev.Floats)
}

// EmbeddingUsage represents the token and cost usage of an embeddings request.
type EmbeddingUsage struct {
	PromptTokens        int                           `json:"prompt_tokens"`
	TotalTokens         int                           `json:"total_tokens"`
	Cost                *float64                      `json:"cost,omitempty"`
	CostDetails         *CostDetails                  `json:"cost_details,omitempty"`
	PromptTokensDetails *EmbeddingPromptTokensDetails `json:"prompt_tokens_details,omitempty"`
	IsBYOK              bool                          `json:"is_byok"`
}

// EmbeddingPromptTokensDetails represents the per-modality token breakdown of
// an embeddings request. It is only present when the input contains two or more
// modalities and the upstream provider returns modality-level usage data. Only
// non-zero modality counts are included.
type EmbeddingPromptTokensDetails struct {
	AudioTokens *int `json:"audio_tokens,omitempty"`
	FileTokens  *int `json:"file_tokens,omitempty"`
	ImageTokens *int `json:"image_tokens,omitempty"`
	TextTokens  *int `json:"text_tokens,omitempty"`
	VideoTokens *int `json:"video_tokens,omitempty"`
}

// EmbeddingEncodingFormat is the encoding of the returned embedding vectors.
type EmbeddingEncodingFormat string

const (
	EmbeddingEncodingFormatFloat  EmbeddingEncodingFormat = "float"
	EmbeddingEncodingFormatBase64 EmbeddingEncodingFormat = "base64"
)

// EmbeddingObject is the object type of an embeddings response or embedding.
type EmbeddingObject string

const (
	EmbeddingObjectList      EmbeddingObject = "list"
	EmbeddingObjectEmbedding EmbeddingObject = "embedding"
)

// EmbeddingContentPartType is the type of a multimodal embedding content part.
type EmbeddingContentPartType string

const (
	EmbeddingContentPartTypeText       EmbeddingContentPartType = "text"
	EmbeddingContentPartTypeImageURL   EmbeddingContentPartType = "image_url"
	EmbeddingContentPartTypeInputAudio EmbeddingContentPartType = "input_audio"
	EmbeddingContentPartTypeInputVideo EmbeddingContentPartType = "input_video"
	EmbeddingContentPartTypeInputFile  EmbeddingContentPartType = "input_file"
)
