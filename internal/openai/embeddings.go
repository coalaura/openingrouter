package openai

// EmbeddingObject is the object type of an embedding.
type EmbeddingObject string

const (
	EmbeddingObjectEmbedding EmbeddingObject = "embedding"
)

// EmbeddingEncodingFormat is the encoding of the returned embedding vectors.
type EmbeddingEncodingFormat string

const (
	EmbeddingEncodingFormatFloat  EmbeddingEncodingFormat = "float"
	EmbeddingEncodingFormatBase64 EmbeddingEncodingFormat = "base64"
)

// EmbeddingInput represents text, token, or batch input(s) to embed. Exactly one
// of the fields should be set. TokenArrays takes precedence over Tokens, Texts
// and Text when more than one is set.
type EmbeddingInput struct {
	Text        string
	Texts       []string
	Tokens      []int
	TokenArrays [][]int
}

// MarshalJSON implements the json.Marshaler interface for EmbeddingInput.
func (ei EmbeddingInput) MarshalJSON() ([]byte, error) {
	switch {
	case len(ei.TokenArrays) > 0:
		return jsonMarshal(ei.TokenArrays)
	case len(ei.Tokens) > 0:
		return jsonMarshal(ei.Tokens)
	case len(ei.Texts) > 0:
		return jsonMarshal(ei.Texts)
	default:
		return jsonMarshal(ei.Text)
	}
}

// EmbeddingRequest is the request body of the embeddings endpoint. Model and
// Input are required, zero values and nil pointers of the remaining fields are
// omitted.
type EmbeddingRequest struct {
	Model string `json:"model"`
	Input EmbeddingInput `json:"input"`

	Dimensions     *int                    `json:"dimensions,omitempty"`
	EncodingFormat EmbeddingEncodingFormat `json:"encoding_format,omitempty"`
	User           string                  `json:"user,omitempty"`
}

// EmbeddingResponse is the root response of the embeddings endpoint.
type EmbeddingResponse struct {
	ID     string     `json:"id,omitempty"`
	Object string     `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string     `json:"model"`
	Usage  *EmbeddingUsage `json:"usage,omitempty"`
}

// Embedding represents a single embedding object of an embeddings response.
type Embedding struct {
	Object    EmbeddingObject `json:"object"`
	Embedding []float64       `json:"embedding"`
	Index     int             `json:"index"`
}

// EmbeddingUsage represents the token usage of an embeddings request.
type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}
