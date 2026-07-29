package openingrouter

import "strings"

// ImageGenerationRequest represents the request body of the image generation
// endpoint. Model and Prompt are required, zero values and nil pointers of the
// remaining fields are omitted from the request. Size is a shorthand for the
// output dimensions: a tier ("2K", "4K") is equivalent to Resolution and
// combines with AspectRatio, an explicit pixel size ("2048x2048") is
// authoritative and is rejected alongside a mismatched Resolution or
// AspectRatio. At most 16 InputReferences are accepted.
type ImageGenerationRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`

	AspectRatio       ImageAspectRatio     `json:"aspect_ratio,omitempty"`
	Background        ImageBackground      `json:"background,omitempty"`
	InputReferences   []ContentPartImage   `json:"input_references,omitempty"`
	Count             *int                 `json:"n,omitempty"`
	OutputCompression *int                 `json:"output_compression,omitempty"`
	OutputFormat      ImageOutputFormat    `json:"output_format,omitempty"`
	Provider          *ProviderPreferences `json:"provider,omitempty"`
	Quality           ImageQuality         `json:"quality,omitempty"`
	Resolution        ImageResolution      `json:"resolution,omitempty"`
	Seed              *int                 `json:"seed,omitempty"`
	Size              string               `json:"size,omitempty"`
	Stream            *bool                `json:"stream,omitempty"`
}

// ImageGenerationResponse is the root response for the image generation endpoint.
type ImageGenerationResponse struct {
	Created int64            `json:"created"`
	Data    []GeneratedImage `json:"data"`
	Usage   *Usage           `json:"usage,omitempty"`
}

// GeneratedImage represents a single generated image. MediaType is omitted if
// the format could not be determined. For svg output the markup is utf-8
// encoded inside B64JSON.
type GeneratedImage struct {
	B64JSON   string `json:"b64_json"`
	MediaType string `json:"media_type,omitempty"`
}

// ImageStreamEvent represents a single event of a streaming image generation
// request. Type determines which of the remaining fields are populated:
// B64JSON and PartialImageIndex for partial images, Text and Phase for text
// chunks, B64JSON, MediaType, Created and Usage for the completed event and
// Error for the error event.
type ImageStreamEvent struct {
	Type              ImageStreamEventType `json:"type"`
	B64JSON           string               `json:"b64_json,omitempty"`
	MediaType         string               `json:"media_type,omitempty"`
	PartialImageIndex int                  `json:"partial_image_index,omitempty"`
	Text              string               `json:"text,omitempty"`
	Phase             ImageStreamPhase     `json:"phase,omitempty"`
	Created           int64                `json:"created,omitempty"`
	Usage             *Usage               `json:"usage,omitempty"`
	Error             *ImageStreamError    `json:"error,omitempty"`
}

func (e ImageStreamEvent) streamError() error {
	if e.Error == nil {
		return nil
	}

	return e.Error
}

// ImageStreamError represents the provider error details of a streaming
// generation that failed after the response started.
type ImageStreamError struct {
	Message string  `json:"message"`
	Code    *string `json:"code"`
	Param   *string `json:"param"`
	Type    *string `json:"type"`
}

// Error returns the formatted string representation of the streaming image error.
func (e *ImageStreamError) Error() string {
	message := parseSubErrorMessage(e.Message)

	if e.Code != nil && *e.Code != "" {
		var sb strings.Builder

		sb.Grow(16 + len(*e.Code) + 2 + len(message))

		sb.WriteString("openrouter code ")
		sb.WriteString(*e.Code)
		sb.WriteString(": ")
		sb.WriteString(message)

		return sb.String()
	}

	var sb strings.Builder

	sb.Grow(12 + len(message))

	sb.WriteString("openrouter: ")
	sb.WriteString(message)

	return sb.String()
}

// ImageAspectRatio is a normalized aspect ratio of a generated image. Providers
// clamp to their supported subset.
type ImageAspectRatio string

const (
	ImageAspectRatio1x1        ImageAspectRatio = "1:1"
	ImageAspectRatio1x2        ImageAspectRatio = "1:2"
	ImageAspectRatio1x4        ImageAspectRatio = "1:4"
	ImageAspectRatio1x8        ImageAspectRatio = "1:8"
	ImageAspectRatio2x1        ImageAspectRatio = "2:1"
	ImageAspectRatio2x3        ImageAspectRatio = "2:3"
	ImageAspectRatio3x2        ImageAspectRatio = "3:2"
	ImageAspectRatio3x4        ImageAspectRatio = "3:4"
	ImageAspectRatio4x1        ImageAspectRatio = "4:1"
	ImageAspectRatio4x3        ImageAspectRatio = "4:3"
	ImageAspectRatio4x5        ImageAspectRatio = "4:5"
	ImageAspectRatio5x4        ImageAspectRatio = "5:4"
	ImageAspectRatio8x1        ImageAspectRatio = "8:1"
	ImageAspectRatio9x16       ImageAspectRatio = "9:16"
	ImageAspectRatio16x9       ImageAspectRatio = "16:9"
	ImageAspectRatio9x19Point5 ImageAspectRatio = "9:19.5"
	ImageAspectRatio19Point5x9 ImageAspectRatio = "19.5:9"
	ImageAspectRatio9x20       ImageAspectRatio = "9:20"
	ImageAspectRatio20x9       ImageAspectRatio = "20:9"
	ImageAspectRatio9x21       ImageAspectRatio = "9:21"
	ImageAspectRatio21x9       ImageAspectRatio = "21:9"
	ImageAspectRatioAuto       ImageAspectRatio = "auto"
)

// ImageBackground is the background treatment of a generated image. Transparent
// requires an output format that supports alpha (png or webp).
type ImageBackground string

const (
	ImageBackgroundAuto        ImageBackground = "auto"
	ImageBackgroundTransparent ImageBackground = "transparent"
	ImageBackgroundOpaque      ImageBackground = "opaque"
)

// ImageOutputFormat is the encoding of the returned image bytes. Svg is
// supported by vectorization models only.
type ImageOutputFormat string

const (
	ImageOutputFormatPNG  ImageOutputFormat = "png"
	ImageOutputFormatJPEG ImageOutputFormat = "jpeg"
	ImageOutputFormatWebP ImageOutputFormat = "webp"
	ImageOutputFormatSVG  ImageOutputFormat = "svg"
)

// ImageQuality is the rendering quality of a generated image. Providers without
// a quality knob ignore it.
type ImageQuality string

const (
	ImageQualityAuto   ImageQuality = "auto"
	ImageQualityLow    ImageQuality = "low"
	ImageQualityMedium ImageQuality = "medium"
	ImageQualityHigh   ImageQuality = "high"
)

// ImageResolution is a normalized resolution tier of a generated image. The
// concrete pixel dimensions are derived per provider.
type ImageResolution string

const (
	ImageResolution512 ImageResolution = "512"
	ImageResolution1K  ImageResolution = "1K"
	ImageResolution2K  ImageResolution = "2K"
	ImageResolution4K  ImageResolution = "4K"
)

// ImageStreamEventType is the type of a streaming image generation event.
type ImageStreamEventType string

const (
	ImageStreamEventTypePartialImage ImageStreamEventType = "image_generation.partial_image"
	ImageStreamEventTypeTextChunk    ImageStreamEventType = "image_generation.text_chunk"
	ImageStreamEventTypeCompleted    ImageStreamEventType = "image_generation.completed"
	ImageStreamEventTypeError        ImageStreamEventType = "error"
)

// ImageStreamPhase is the generation phase a text chunk belongs to. Content is
// the renderable output, reasoning and draft are intermediate provider phases.
type ImageStreamPhase string

const (
	ImageStreamPhaseContent   ImageStreamPhase = "content"
	ImageStreamPhaseReasoning ImageStreamPhase = "reasoning"
	ImageStreamPhaseDraft     ImageStreamPhase = "draft"
)
