package openingrouter

// ListEmbeddingModelsOptions holds the optional query parameters of the
// embeddings models list endpoint. Zero values and nil pointers are omitted
// from the request. When both Offset and Limit are omitted, the full list is
// returned.
type ListEmbeddingModelsOptions struct {
	Offset *int `url:"offset,omitempty"`
	Limit  *int `url:"limit,omitempty"`
}
