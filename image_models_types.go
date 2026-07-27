package openingrouter

// ImageModelsListResponse is the root response for the image models list endpoint.
type ImageModelsListResponse struct {
	Data []ImageModel `json:"data"`
}

// ImageModel represents an image generation model.
type ImageModel struct {
	ID                  string                   `json:"id"`
	Name                string                   `json:"name"`
	Description         string                   `json:"description"`
	Created             int64                    `json:"created"`
	Architecture        ImageModelArchitecture   `json:"architecture"`
	SupportedParameters ImageSupportedParameters `json:"supported_parameters"`
	SupportsStreaming   bool                     `json:"supports_streaming"`
	Endpoints           string                   `json:"endpoints"`
}

// ImageModelArchitecture represents the architecture information of an image
// generation model.
type ImageModelArchitecture struct {
	InputModalities  []InputModality  `json:"input_modalities"`
	OutputModalities []OutputModality `json:"output_modalities"`
}

// ImageSupportedParameters represents the union of supported parameters across
// every endpoint of an image generation model, keyed by parameter name. It is a
// coarse discovery aid, the definitive per-endpoint set is behind the endpoints
// url of the model.
type ImageSupportedParameters map[string]ImageCapability

// ImageCapability is a typed descriptor for one supported request parameter. Type
// determines which of the remaining fields are populated: Values for enum, Min
// and Max for range, none for boolean.
type ImageCapability struct {
	Type   ImageCapabilityType `json:"type"`
	Values []string            `json:"values,omitempty"`
	Min    *float64            `json:"min,omitempty"`
	Max    *float64            `json:"max,omitempty"`
}

// ImageCapabilityType is the kind of a supported parameter descriptor.
type ImageCapabilityType string

const (
	ImageCapabilityTypeBoolean ImageCapabilityType = "boolean"
	ImageCapabilityTypeEnum    ImageCapabilityType = "enum"
	ImageCapabilityTypeRange   ImageCapabilityType = "range"
)
