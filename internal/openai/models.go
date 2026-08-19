// Package openai holds the wire types of the OpenAI api. They are the raw
// request and response shapes of an OpenAI compatible endpoint and are kept
// separate from the openingrouter types, which the public api exposes.
package openai

// ModelObject is the object type of a model, always "model".
type ModelObject string

const (
	ModelObjectModel ModelObject = "model"
)

// ListObject is the object type of a list response, always "list".
type ListObject string

const (
	ListObjectList ListObject = "list"
)

// Model represents a single model offering of an OpenAI compatible endpoint.
type Model struct {
	ID           string  `json:"id"`
	Created      int64   `json:"created"`
	Object       string  `json:"object"`
	OwnedBy      string  `json:"owned_by"`
	ShutdownDate *string `json:"shutdown_date,omitempty"`
}

// ModelsList is the root response of the list models endpoint.
type ModelsList struct {
	Object string  `json:"object"`
	Models []Model `json:"models"`
}
