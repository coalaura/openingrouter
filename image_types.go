package openingrouter

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

	AspectRatio       ImageAspectRatio                    `json:"aspect_ratio,omitempty"`
	Background        ImageBackground                     `json:"background,omitempty"`
	InputReferences   []ContentPartImage                  `json:"input_references,omitempty"`
	Count             *int                                `json:"n,omitempty"`
	OutputCompression *int                                `json:"output_compression,omitempty"`
	OutputFormat      ImageOutputFormat                   `json:"output_format,omitempty"`
	Provider          *ImageGenerationProviderPreferences `json:"provider,omitempty"`
	Quality           ImageQuality                        `json:"quality,omitempty"`
	Resolution        ImageResolution                     `json:"resolution,omitempty"`
	Seed              *int                                `json:"seed,omitempty"`
	Size              string                              `json:"size,omitempty"`
	Stream            *bool                               `json:"stream,omitempty"`
}

// ContentPartImage represents a reference image guiding image-to-image
// generation, as a base64 data url or an HTTP(S) url.
type ContentPartImage struct {
	Type     ContentPartType     `json:"type"`
	ImageURL ContentPartImageURL `json:"image_url"`
}

// ContentPartImageURL holds the url of a ContentPartImage.
type ContentPartImageURL struct {
	URL string `json:"url"`
}

// ImageGenerationProviderPreferences represents the provider routing preferences
// and provider specific passthrough configuration of an image generation
// request. Order, Only and Ignore are merged with the account-wide settings.
type ImageGenerationProviderPreferences struct {
	Order          []string            `json:"order,omitempty"`
	Only           []string            `json:"only,omitempty"`
	Ignore         []string            `json:"ignore,omitempty"`
	AllowFallbacks *bool               `json:"allow_fallbacks,omitempty"`
	Sort           *ProviderSortConfig `json:"sort,omitempty"`
	Options        ProviderOptions     `json:"options,omitempty"`
}

// ProviderSortConfig represents the sorting strategy used for a request when no
// explicit order is given. Setting it disables load balancing.
type ProviderSortConfig struct {
	By        ProviderSort          `json:"by,omitempty"`
	Partition ProviderSortPartition `json:"partition,omitempty"`
}

// ProviderOptions holds provider specific options keyed by provider slug. Only
// the options of the matched provider are forwarded, the rest are ignored and
// unrecognized keys are silently dropped.
type ProviderOptions map[string]map[string]any

// ImageGenerationResponse is the root response for the image generation endpoint.
type ImageGenerationResponse struct {
	Created int64                 `json:"created"`
	Data    []GeneratedImage      `json:"data"`
	Usage   *ImageGenerationUsage `json:"usage,omitempty"`
}

// GeneratedImage represents a single generated image. MediaType is omitted if
// the format could not be determined. For svg output the markup is utf-8
// encoded inside B64JSON.
type GeneratedImage struct {
	B64JSON   string `json:"b64_json"`
	MediaType string `json:"media_type,omitempty"`
}

// ImageGenerationUsage represents the token and cost usage of an image
// generation request, when available.
type ImageGenerationUsage struct {
	PromptTokens            int                       `json:"prompt_tokens"`
	CompletionTokens        int                       `json:"completion_tokens"`
	TotalTokens             int                       `json:"total_tokens"`
	Cost                    *float64                  `json:"cost"`
	CostDetails             *CostDetails              `json:"cost_details"`
	PromptTokensDetails     *PromptTokensDetails      `json:"prompt_tokens_details"`
	CompletionTokensDetails *CompletionTokensDetails  `json:"completion_tokens_details"`
	ServerToolUse           *ServerToolUse            `json:"server_tool_use"`
	CacheCreation           *AnthropicCacheCreation   `json:"cache_creation"`
	Iterations              []AnthropicUsageIteration `json:"iterations"`
	Speed                   AnthropicSpeed            `json:"speed,omitempty"`
	ServiceTier             *string                   `json:"service_tier"`
	IsBYOK                  bool                      `json:"is_byok"`
}

// CostDetails represents the breakdown of upstream inference costs.
type CostDetails struct {
	UpstreamInferenceCost            *float64 `json:"upstream_inference_cost"`
	UpstreamInferencePromptCost      float64  `json:"upstream_inference_prompt_cost"`
	UpstreamInferenceCompletionsCost float64  `json:"upstream_inference_completions_cost"`
}

// PromptTokensDetails represents the breakdown of tokens used in the prompt.
// CacheWriteTokens is only returned for models with explicit caching and cache
// write pricing.
type PromptTokensDetails struct {
	CachedTokens     *int `json:"cached_tokens"`
	CacheWriteTokens *int `json:"cache_write_tokens"`
	AudioTokens      *int `json:"audio_tokens"`
	FileTokens       *int `json:"file_tokens"`
	VideoTokens      *int `json:"video_tokens"`
}

// CompletionTokensDetails represents the breakdown of tokens generated by the
// model.
type CompletionTokensDetails struct {
	ReasoningTokens *int `json:"reasoning_tokens"`
	AudioTokens     *int `json:"audio_tokens"`
	ImageTokens     *int `json:"image_tokens"`
}

// ServerToolUse represents the usage of server-side tool execution. A
// server-orchestrated web search is counted in both ToolCallsRequested and
// WebSearchRequests, provider-native web search may report WebSearchRequests
// only, so the two must not be summed.
type ServerToolUse struct {
	ToolCallsRequested *int `json:"tool_calls_requested"`
	ToolCallsExecuted  *int `json:"tool_calls_executed"`
	WebSearchRequests  *int `json:"web_search_requests"`
}

// AnthropicCacheCreation represents the cache write tokens of a request, split
// by cache ttl.
type AnthropicCacheCreation struct {
	Ephemeral5MInputTokens int `json:"ephemeral_5m_input_tokens"`
	Ephemeral1HInputTokens int `json:"ephemeral_1h_input_tokens"`
}

// AnthropicIterationCacheCreation represents the cache write tokens of a single
// usage iteration, split by cache ttl.
type AnthropicIterationCacheCreation struct {
	Ephemeral5MInputTokens int `json:"ephemeral_5m_input_tokens"`
	Ephemeral1HInputTokens int `json:"ephemeral_1h_input_tokens"`
}

// AnthropicUsageIteration represents the usage of a single iteration of a
// request. Model is only populated for message and advisor message iterations,
// unknown iteration types are passed through as-is.
type AnthropicUsageIteration struct {
	Type                     AnthropicUsageIterationType      `json:"type"`
	Model                    string                           `json:"model,omitempty"`
	InputTokens              int                              `json:"input_tokens"`
	OutputTokens             int                              `json:"output_tokens"`
	CacheCreationInputTokens int                              `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int                              `json:"cache_read_input_tokens"`
	CacheCreation            *AnthropicIterationCacheCreation `json:"cache_creation"`
}

// ImageStreamEvent represents a single event of a streaming image generation
// request. Type determines which of the remaining fields are populated:
// B64JSON and PartialImageIndex for partial images, Text and Phase for text
// chunks, B64JSON, MediaType, Created and Usage for the completed event and
// Error for the error event.
type ImageStreamEvent struct {
	Type              ImageStreamEventType  `json:"type"`
	B64JSON           string                `json:"b64_json,omitempty"`
	MediaType         string                `json:"media_type,omitempty"`
	PartialImageIndex int                   `json:"partial_image_index,omitempty"`
	Text              string                `json:"text,omitempty"`
	Phase             ImageStreamPhase      `json:"phase,omitempty"`
	Created           int64                 `json:"created,omitempty"`
	Usage             *ImageGenerationUsage `json:"usage,omitempty"`
	Error             *ImageStreamError     `json:"error,omitempty"`
}

// ImageStreamError represents the provider error details of a streaming
// generation that failed after the response started.
type ImageStreamError struct {
	Message string  `json:"message"`
	Code    *string `json:"code"`
	Param   *string `json:"param"`
	Type    *string `json:"type"`
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

// ContentPartType is the type of a content part.
type ContentPartType string

const (
	ContentPartTypeImageURL ContentPartType = "image_url"
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

// AnthropicSpeed is the speed tier a request was served with.
type AnthropicSpeed string

const (
	AnthropicSpeedFast     AnthropicSpeed = "fast"
	AnthropicSpeedStandard AnthropicSpeed = "standard"
)

// AnthropicUsageIterationType is the type of a single usage iteration.
type AnthropicUsageIterationType string

const (
	AnthropicUsageIterationTypeCompaction     AnthropicUsageIterationType = "compaction"
	AnthropicUsageIterationTypeMessage        AnthropicUsageIterationType = "message"
	AnthropicUsageIterationTypeAdvisorMessage AnthropicUsageIterationType = "advisor_message"
)

// ProviderSort is the sorting strategy used to pick an endpoint.
type ProviderSort string

const (
	ProviderSortPrice      ProviderSort = "price"
	ProviderSortThroughput ProviderSort = "throughput"
	ProviderSortLatency    ProviderSort = "latency"
	ProviderSortExacto     ProviderSort = "exacto"
)

// ProviderSortPartition is the partitioning strategy applied before sorting.
// Model groups endpoints by model so fallback models remain fallbacks, none
// sorts all endpoints together regardless of model.
type ProviderSortPartition string

const (
	ProviderSortPartitionModel ProviderSortPartition = "model"
	ProviderSortPartitionNone  ProviderSortPartition = "none"
)
