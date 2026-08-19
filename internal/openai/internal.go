package openai

import "encoding/json"

// jsonMarshal marshals v to json, a small indirection so the wire types can
// share one encoding path.
func jsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
