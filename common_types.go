package openingrouter

type OpenRouterResponse[T any] struct {
	Data T `json:"data"`
}
