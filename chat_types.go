package openingrouter

import (
	"bytes"
	"encoding/json"
	"strings"
)

// ChatCompletionRequest represents the request body of the chat completion
// endpoint. Messages is required, Model may be omitted when a router plugin or
// preset decides the route. Zero values and nil pointers of the remaining
// fields are omitted from the request. MetadataLevel is sent as the
// X-OpenRouter-Metadata header and is not part of the body.
type ChatCompletionRequest struct {
	Model    string        `json:"model,omitempty"`
	Messages []ChatMessage `json:"messages"`

	CacheControl        *AnthropicCacheControl  `json:"cache_control,omitempty"`
	Debug               *ChatDebugOptions       `json:"debug,omitempty"`
	FrequencyPenalty    *float64                `json:"frequency_penalty,omitempty"`
	ImageConfig         ChatImageConfig         `json:"image_config,omitempty"`
	LogitBias           map[string]float64      `json:"logit_bias,omitempty"`
	Logprobs            *bool                   `json:"logprobs,omitempty"`
	MaxCompletionTokens *int                    `json:"max_completion_tokens,omitempty"`
	MaxTokens           *int                    `json:"max_tokens,omitempty"`
	Metadata            map[string]string       `json:"metadata,omitempty"`
	MinP                *float64                `json:"min_p,omitempty"`
	Modalities          []OutputModality        `json:"modalities,omitempty"`
	Models              []string                `json:"models,omitempty"`
	ParallelToolCalls   *bool                   `json:"parallel_tool_calls,omitempty"`
	Plugins             []ChatPlugin            `json:"plugins,omitempty"`
	Prediction          *ChatPrediction         `json:"prediction,omitempty"`
	PresencePenalty     *float64                `json:"presence_penalty,omitempty"`
	PromptCacheKey      string                  `json:"prompt_cache_key,omitempty"`
	PromptCacheOptions  *ChatPromptCacheOptions `json:"prompt_cache_options,omitempty"`
	Provider            *ProviderPreferences    `json:"provider,omitempty"`
	Reasoning           *ChatReasoningConfig    `json:"reasoning,omitempty"`
	ReasoningEffort     ReasoningEffort         `json:"reasoning_effort,omitempty"`
	RepetitionPenalty   *float64                `json:"repetition_penalty,omitempty"`
	ResponseFormat      *ChatResponseFormat     `json:"response_format,omitempty"`
	Route               ChatRoute               `json:"route,omitempty"`
	Seed                *int                    `json:"seed,omitempty"`
	ServiceTier         ChatServiceTier         `json:"service_tier,omitempty"`
	SessionID           string                  `json:"session_id,omitempty"`
	Stop                []string                `json:"stop,omitempty"`
	StopServerToolsWhen []ChatStopCondition     `json:"stop_server_tools_when,omitempty"`
	Stream              *bool                   `json:"stream,omitempty"`
	StreamOptions       *ChatStreamOptions      `json:"stream_options,omitempty"`
	Temperature         *float64                `json:"temperature,omitempty"`
	ToolChoice          *ChatToolChoice         `json:"tool_choice,omitempty"`
	Tools               []ChatTool              `json:"tools,omitempty"`
	TopA                *float64                `json:"top_a,omitempty"`
	TopK                *int                    `json:"top_k,omitempty"`
	TopLogprobs         *int                    `json:"top_logprobs,omitempty"`
	TopP                *float64                `json:"top_p,omitempty"`
	Trace               *ChatTraceConfig        `json:"trace,omitempty"`
	User                string                  `json:"user,omitempty"`

	MetadataLevel ChatMetadataLevel `json:"-"`
}

// ChatDebugOptions holds the debug options of a request. They only take effect
// on streaming requests.
type ChatDebugOptions struct {
	EchoUpstreamBody bool `json:"echo_upstream_body,omitempty"`
}

// ChatImageConfig holds provider specific image generation options keyed by
// option name (aspect_ratio, quality, size, ...). Unrecognized keys are
// forwarded as-is and ignored by providers that do not support them.
type ChatImageConfig map[string]any

// ChatStreamOptions holds the streaming options of a request.
type ChatStreamOptions struct {
	// Deprecated: this field has no effect, full usage details are always
	// included.
	IncludeUsage *bool `json:"include_usage,omitempty"`
}

// ChatReasoningConfig represents the reasoning configuration of a request.
// Effort cannot differ from the shorthand ChatCompletionRequest.ReasoningEffort.
type ChatReasoningConfig struct {
	Effort  ReasoningEffort      `json:"effort,omitempty"`
	Summary ChatReasoningSummary `json:"summary,omitempty"`
}

// ChatPrediction represents static predicted output content. Supported models
// use it to reduce latency when much of the response is known in advance.
type ChatPrediction struct {
	Type    ChatPredictionType `json:"type"`
	Content ChatContent        `json:"content"`
}

// ChatPromptCacheOptions represents the request level prompt cache controls.
// ChatPromptCacheModeExplicit disables provider managed breakpoints so only
// blocks marked with a ChatPromptCacheBreakpoint are cached. Only supported by
// OpenAI GPT-5.6 and newer.
type ChatPromptCacheOptions struct {
	Mode ChatPromptCacheMode `json:"mode"`
	TTL  string              `json:"ttl,omitempty"`
}

// ChatPromptCacheBreakpoint marks an explicit prompt cache boundary on a
// content block. Everything through the block carrying the marker is part of
// the candidate cached prefix. It is interchangeable with AnthropicCacheControl,
// OpenRouter converts between the two based on the serving provider.
type ChatPromptCacheBreakpoint struct {
	Mode ChatPromptCacheMode `json:"mode"`
}

// AnthropicCacheControl enables automatic prompt caching. At the top level of a
// request the last cacheable block is used as the breakpoint, on a content block
// it marks an explicit breakpoint.
type AnthropicCacheControl struct {
	Type AnthropicCacheControlType `json:"type"`
	TTL  AnthropicCacheTTL         `json:"ttl,omitempty"`
}

// ChatTraceConfig holds the observability metadata of a request. The known keys
// receive special handling, they are forwarded to the configured broadcast
// destinations.
type ChatTraceConfig struct {
	TraceID        string `json:"trace_id,omitempty"`
	TraceName      string `json:"trace_name,omitempty"`
	SpanName       string `json:"span_name,omitempty"`
	GenerationName string `json:"generation_name,omitempty"`
	ParentSpanID   string `json:"parent_span_id,omitempty"`
}

// ChatResponseFormat represents the response format configuration of a request.
// Type determines which of the remaining fields are used: JSONSchema for
// json_schema, Grammar for grammar, none for text, json_object and python.
type ChatResponseFormat struct {
	Type       ChatResponseFormatType `json:"type"`
	JSONSchema *ChatJSONSchema        `json:"json_schema,omitempty"`
	Grammar    string                 `json:"grammar,omitempty"`
}

// ChatJSONSchema represents the schema of a json_schema response format. Name is
// required and limited to 64 characters of a-z, A-Z, 0-9, underscores and dashes.
type ChatJSONSchema struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Schema      map[string]any `json:"schema,omitempty"`
	Strict      *bool          `json:"strict,omitempty"`
}

// ChatStopCondition represents a single stop condition of the server tool agent
// loop. Any condition firing halts the loop. Type determines which of the
// remaining fields are used.
type ChatStopCondition struct {
	Type             ChatStopConditionType `json:"type"`
	StepCount        *int                  `json:"step_count,omitempty"`
	ToolName         string                `json:"tool_name,omitempty"`
	MaxTokens        *int                  `json:"max_tokens,omitempty"`
	MaxCostInDollars *float64              `json:"max_cost_in_dollars,omitempty"`
	Reason           string                `json:"reason,omitempty"`
}

// ChatMessage represents a single message of a conversation, for both requests
// and responses. Role determines which of the remaining fields are meaningful:
// Content and Name for system, developer and user messages, Content and
// ToolCallID for tool messages, everything else for assistant messages.
type ChatMessage struct {
	Role    ChatRole    `json:"role"`
	Content ChatContent `json:"content,omitzero"`
	Name    string      `json:"name,omitempty"`

	Audio            *ChatAudioOutput      `json:"audio,omitempty"`
	Images           []ChatAssistantImage  `json:"images,omitempty"`
	Model            string                `json:"model,omitempty"`
	Reasoning        string                `json:"reasoning,omitempty"`
	ReasoningDetails []ChatReasoningDetail `json:"reasoning_details,omitempty"`
	Refusal          string                `json:"refusal,omitempty"`
	ToolCalls        []ChatToolCall        `json:"tool_calls,omitempty"`

	ToolCallID string `json:"tool_call_id,omitempty"`
}

// ChatContent represents the content of a message, encoded either as a plain
// string or as a list of content parts. Parts takes precedence over Text when
// both are set.
type ChatContent struct {
	Text  string
	Parts []ChatContentPart
}

// MarshalJSON implements the json.Marshaler interface for ChatContent.
func (cc ChatContent) MarshalJSON() ([]byte, error) {
	if len(cc.Parts) > 0 {
		return json.Marshal(cc.Parts)
	}

	return json.Marshal(cc.Text)
}

// UnmarshalJSON implements the json.Unmarshaler interface for ChatContent.
func (cc *ChatContent) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		*cc = ChatContent{}

		return nil
	}

	if len(data) > 0 && data[0] == '[' {
		return json.Unmarshal(data, &cc.Parts)
	}

	return json.Unmarshal(data, &cc.Text)
}

// String returns the text of the content, joining every text part if it is
// encoded as content parts. Non-text parts are skipped.
func (cc ChatContent) String() string {
	if len(cc.Parts) == 0 {
		return cc.Text
	}

	var size int

	for _, part := range cc.Parts {
		if part.Type == ChatContentPartTypeText {
			size += len(part.Text)
		}
	}

	var sb strings.Builder

	sb.Grow(size)

	for _, part := range cc.Parts {
		if part.Type != ChatContentPartTypeText {
			continue
		}

		sb.WriteString(part.Text)
	}

	return sb.String()
}

// ChatContentPart represents a single content part of a message. Type
// determines which of the remaining fields are used.
type ChatContentPart struct {
	Type                  ChatContentPartType        `json:"type"`
	Text                  string                     `json:"text,omitempty"`
	ImageURL              *ChatContentImageURL       `json:"image_url,omitempty"`
	InputAudio            *ChatContentInputAudio     `json:"input_audio,omitempty"`
	VideoURL              *ChatContentVideoURL       `json:"video_url,omitempty"`
	File                  *ChatContentFile           `json:"file,omitempty"`
	CacheControl          *AnthropicCacheControl     `json:"cache_control,omitempty"`
	PromptCacheBreakpoint *ChatPromptCacheBreakpoint `json:"prompt_cache_breakpoint,omitempty"`
}

// ChatContentImageURL holds the url of an image content part, as a base64 data
// url or an HTTP(S) url.
type ChatContentImageURL struct {
	URL    string          `json:"url"`
	Detail ChatImageDetail `json:"detail,omitempty"`
}

// ChatContentInputAudio holds the base64 encoded audio of an audio content part.
// The supported formats (wav, mp3, flac, m4a, ogg, aiff, aac, pcm16, pcm24)
// vary by provider.
type ChatContentInputAudio struct {
	Data   string `json:"data"`
	Format string `json:"format"`
}

// ChatContentVideoURL holds the url of a video content part, as a base64 data
// url or an HTTP(S) url.
type ChatContentVideoURL struct {
	URL string `json:"url"`
}

// ChatContentFile holds the document of a file content part, either inline as a
// base64 data url or url, or by the id of a previously uploaded file.
type ChatContentFile struct {
	FileData string `json:"file_data,omitempty"`
	FileID   string `json:"file_id,omitempty"`
	Filename string `json:"filename,omitempty"`
}

// ChatAssistantImage represents a single image generated by an image generation
// model.
type ChatAssistantImage struct {
	ImageURL ContentPartImageURL `json:"image_url"`
}

// ChatAudioOutput represents the audio output of an assistant message.
type ChatAudioOutput struct {
	ID         string `json:"id,omitempty"`
	Data       string `json:"data,omitempty"`
	Transcript string `json:"transcript,omitempty"`
	ExpiresAt  int64  `json:"expires_at,omitempty"`
}

// ChatReasoningDetail represents a single reasoning detail of an extended
// thinking model. Type determines which of the remaining fields are populated:
// Summary for summaries, Data for encrypted reasoning, Text and Signature for
// text and ToolName, Arguments, Result and ToolCallID for server tool calls.
type ChatReasoningDetail struct {
	Type       ChatReasoningDetailType `json:"type"`
	ID         string                  `json:"id,omitempty"`
	Index      int                     `json:"index,omitempty"`
	Format     ChatReasoningFormat     `json:"format,omitempty"`
	Summary    string                  `json:"summary,omitempty"`
	Data       string                  `json:"data,omitempty"`
	Text       string                  `json:"text,omitempty"`
	Signature  string                  `json:"signature,omitempty"`
	ToolName   string                  `json:"tool_name,omitempty"`
	Arguments  string                  `json:"arguments,omitempty"`
	Result     string                  `json:"result,omitempty"`
	ToolCallID string                  `json:"tool_call_id,omitempty"`
}

// ChatToolCall represents a tool call made by the assistant.
type ChatToolCall struct {
	ID       string               `json:"id"`
	Type     ChatToolType         `json:"type"`
	Function ChatToolCallFunction `json:"function"`
}

// ChatToolCallFunction holds the name and json encoded arguments of a tool call.
type ChatToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ChatToolChoice represents the tool choice configuration of a request. Mode
// encodes the plain "none", "auto" and "required" choices and takes precedence,
// otherwise Type and Function name the tool to force.
type ChatToolChoice struct {
	Mode     ChatToolChoiceMode      `json:"-"`
	Type     ChatToolType            `json:"type,omitempty"`
	Function *ChatToolChoiceFunction `json:"function,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface for ChatToolChoice.
func (tc ChatToolChoice) MarshalJSON() ([]byte, error) {
	if tc.Mode != "" {
		return json.Marshal(string(tc.Mode))
	}

	type choice ChatToolChoice

	return json.Marshal(choice(tc))
}

// ChatToolChoiceFunction holds the name of the function a named tool choice
// forces.
type ChatToolChoiceFunction struct {
	Name string `json:"name"`
}

// ChatTool is implemented by every tool definition of a chat completion
// request, both regular function tools and OpenRouter built-in server tools.
type ChatTool interface {
	chatTool()
}

// ChatFunctionTool represents a regular function tool.
type ChatFunctionTool struct {
	Type         ChatToolType           `json:"type"`
	Function     ChatFunction           `json:"function"`
	CacheControl *AnthropicCacheControl `json:"cache_control,omitempty"`
}

func (t ChatFunctionTool) chatTool() {}

// ChatFunction represents the definition of a callable function. Name is
// limited to 64 characters of a-z, A-Z, 0-9, underscores and dashes.
type ChatFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
	Strict      *bool          `json:"strict,omitempty"`
}

// ChatAdvisorTool represents the built-in advisor server tool, consulting a
// higher-intelligence advisor model mid-generation. Include multiple entries to
// offer several named advisors, at most one may omit its name.
type ChatAdvisorTool struct {
	Type       ChatToolType           `json:"type"`
	Parameters *ChatAdvisorToolConfig `json:"parameters,omitempty"`
}

func (t ChatAdvisorTool) chatTool() {}

// ChatAdvisorToolConfig holds the configuration of a single advisor server tool
// entry.
type ChatAdvisorToolConfig struct {
	ForwardTranscript   *bool              `json:"forward_transcript,omitempty"`
	Instructions        string             `json:"instructions,omitempty"`
	MaxCompletionTokens *int               `json:"max_completion_tokens,omitempty"`
	MaxToolCalls        *int               `json:"max_tool_calls,omitempty"`
	Model               string             `json:"model,omitempty"`
	Name                string             `json:"name,omitempty"`
	Reasoning           *ChatToolReasoning `json:"reasoning,omitempty"`
	Stream              *bool              `json:"stream,omitempty"`
	Temperature         *float64           `json:"temperature,omitempty"`
	Tools               []ChatNestedTool   `json:"tools,omitempty"`
}

// ChatBashTool represents the built-in bash server tool, running shell commands
// in a sandboxed container.
type ChatBashTool struct {
	Type       ChatToolType        `json:"type"`
	Parameters *ChatBashToolConfig `json:"parameters,omitempty"`
}

func (t ChatBashTool) chatTool() {}

// ChatBashToolConfig holds the configuration of the bash server tool.
// SleepAfterSeconds is idle based, defaults to 900 and is capped at 2592000.
type ChatBashToolConfig struct {
	Engine            ChatBashEngine       `json:"engine,omitempty"`
	Environment       *ChatBashEnvironment `json:"environment,omitempty"`
	SleepAfterSeconds *int                 `json:"sleep_after_seconds,omitempty"`
}

// ChatBashEnvironment represents the execution environment of the bash server
// tool. ContainerID is only used for container references.
type ChatBashEnvironment struct {
	Type        ChatBashEnvironmentType `json:"type"`
	ContainerID string                  `json:"container_id,omitempty"`
}

// ChatDatetimeTool represents the built-in datetime server tool, returning the
// current date and time.
type ChatDatetimeTool struct {
	Type       ChatToolType            `json:"type"`
	Parameters *ChatDatetimeToolConfig `json:"parameters,omitempty"`
}

func (t ChatDatetimeTool) chatTool() {}

// ChatDatetimeToolConfig holds the configuration of the datetime server tool.
// Timezone is an IANA timezone name and defaults to UTC.
type ChatDatetimeToolConfig struct {
	Timezone string `json:"timezone,omitempty"`
}

// ChatFilesTool represents the built-in files server tool, reading, writing,
// editing and listing workspace files. It requires the x-openrouter-file-ids
// request header.
type ChatFilesTool struct {
	Type       ChatToolType         `json:"type"`
	Parameters *ChatFilesToolConfig `json:"parameters,omitempty"`
}

func (t ChatFilesTool) chatTool() {}

// ChatFilesToolConfig holds the configuration of the files server tool.
type ChatFilesToolConfig struct{}

// ChatFusionTool represents the built-in fusion server tool, fanning the prompt
// out to a panel of analysis models and summarizing their output.
type ChatFusionTool struct {
	Type       ChatToolType          `json:"type"`
	Parameters *ChatFusionToolConfig `json:"parameters,omitempty"`
}

func (t ChatFusionTool) chatTool() {}

// ChatFusionToolConfig holds the configuration of the fusion server tool. At
// most 8 analysis models are accepted, MaxToolCalls is capped at 16.
type ChatFusionToolConfig struct {
	AnalysisModels      []string               `json:"analysis_models,omitempty"`
	CacheControl        *AnthropicCacheControl `json:"cache_control,omitempty"`
	MaxCompletionTokens *int                   `json:"max_completion_tokens,omitempty"`
	MaxToolCalls        *int                   `json:"max_tool_calls,omitempty"`
	Model               string                 `json:"model,omitempty"`
	Reasoning           *ChatToolReasoning     `json:"reasoning,omitempty"`
	Temperature         *float64               `json:"temperature,omitempty"`
	Tools               []ChatNestedTool       `json:"tools,omitempty"`
}

// ChatImageGenerationTool represents the built-in image generation server tool.
type ChatImageGenerationTool struct {
	Type       ChatToolType                  `json:"type"`
	Parameters ChatImageGenerationToolConfig `json:"parameters,omitempty"`
}

func (t ChatImageGenerationTool) chatTool() {}

// ChatImageGenerationToolConfig holds the configuration of the image generation
// server tool. It accepts every ChatImageConfig option plus a "model" key,
// which defaults to openai/gpt-5-image.
type ChatImageGenerationToolConfig map[string]any

// ChatSearchModelsTool represents the built-in experimental search models server
// tool, searching and filtering the models available on OpenRouter.
type ChatSearchModelsTool struct {
	Type       ChatToolType                `json:"type"`
	Parameters *ChatSearchModelsToolConfig `json:"parameters,omitempty"`
}

func (t ChatSearchModelsTool) chatTool() {}

// ChatSearchModelsToolConfig holds the configuration of the search models server
// tool. MaxResults defaults to 5 and is capped at 20.
type ChatSearchModelsToolConfig struct {
	MaxResults *int `json:"max_results,omitempty"`
}

// ChatSubagentTool represents the built-in subagent server tool, delegating
// self-contained tasks to a smaller, cheaper worker model.
type ChatSubagentTool struct {
	Type       ChatToolType            `json:"type"`
	Parameters *ChatSubagentToolConfig `json:"parameters,omitempty"`
}

func (t ChatSubagentTool) chatTool() {}

// ChatSubagentToolConfig holds the configuration of a single subagent server
// tool entry.
type ChatSubagentToolConfig struct {
	Instructions        string             `json:"instructions,omitempty"`
	MaxCompletionTokens *int               `json:"max_completion_tokens,omitempty"`
	MaxToolCalls        *int               `json:"max_tool_calls,omitempty"`
	Model               string             `json:"model,omitempty"`
	Name                string             `json:"name,omitempty"`
	Reasoning           *ChatToolReasoning `json:"reasoning,omitempty"`
	Temperature         *float64           `json:"temperature,omitempty"`
	Tools               []ChatNestedTool   `json:"tools,omitempty"`
}

// ChatWebFetchTool represents the built-in web fetch server tool, fetching the
// full content of a web page or PDF.
type ChatWebFetchTool struct {
	Type       ChatToolType            `json:"type"`
	Parameters *ChatWebFetchToolConfig `json:"parameters,omitempty"`
}

func (t ChatWebFetchTool) chatTool() {}

// ChatWebFetchToolConfig holds the configuration of the web fetch server tool.
type ChatWebFetchToolConfig struct {
	AllowedDomains   []string           `json:"allowed_domains,omitempty"`
	BlockedDomains   []string           `json:"blocked_domains,omitempty"`
	Engine           ChatWebFetchEngine `json:"engine,omitempty"`
	MaxContentTokens *int               `json:"max_content_tokens,omitempty"`
	MaxUses          *int               `json:"max_uses,omitempty"`
}

// ChatWebSearchTool represents the built-in web search server tool.
type ChatWebSearchTool struct {
	Type       ChatToolType             `json:"type"`
	Parameters *ChatWebSearchToolConfig `json:"parameters,omitempty"`
}

func (t ChatWebSearchTool) chatTool() {}

// ChatWebSearchToolConfig holds the configuration of the web search server tool.
// AllowedDomains and ExcludedDomains are mutually exclusive, MaxCharacters takes
// precedence over SearchContextSize.
type ChatWebSearchToolConfig struct {
	AllowedDomains    []string                   `json:"allowed_domains,omitempty"`
	Engine            ChatWebSearchEngine        `json:"engine,omitempty"`
	ExcludedDomains   []string                   `json:"excluded_domains,omitempty"`
	MaxCharacters     *int                       `json:"max_characters,omitempty"`
	MaxResults        *int                       `json:"max_results,omitempty"`
	MaxTotalResults   *int                       `json:"max_total_results,omitempty"`
	MaxUses           *int                       `json:"max_uses,omitempty"`
	SearchContextSize ChatSearchContextSize      `json:"search_context_size,omitempty"`
	UserLocation      *ChatWebSearchUserLocation `json:"user_location,omitempty"`
}

// ChatWebSearchShorthandTool represents a web search tool declared with the
// OpenAI Responses API syntax. It is converted to the built-in web search server
// tool, Parameters overrides the flattened options.
type ChatWebSearchShorthandTool struct {
	Type              ChatToolType               `json:"type"`
	AllowedDomains    []string                   `json:"allowed_domains,omitempty"`
	Engine            ChatWebSearchEngine        `json:"engine,omitempty"`
	ExcludedDomains   []string                   `json:"excluded_domains,omitempty"`
	MaxCharacters     *int                       `json:"max_characters,omitempty"`
	MaxResults        *int                       `json:"max_results,omitempty"`
	MaxTotalResults   *int                       `json:"max_total_results,omitempty"`
	MaxUses           *int                       `json:"max_uses,omitempty"`
	Parameters        *ChatWebSearchToolConfig   `json:"parameters,omitempty"`
	SearchContextSize ChatSearchContextSize      `json:"search_context_size,omitempty"`
	UserLocation      *ChatWebSearchUserLocation `json:"user_location,omitempty"`
}

func (t ChatWebSearchShorthandTool) chatTool() {}

// ChatWebSearchUserLocation represents the approximate user location used to
// bias search results.
type ChatWebSearchUserLocation struct {
	Type     ChatUserLocationType `json:"type"`
	City     string               `json:"city,omitempty"`
	Country  string               `json:"country,omitempty"`
	Region   string               `json:"region,omitempty"`
	Timezone string               `json:"timezone,omitempty"`
}

// ChatToolReasoning represents the reasoning configuration forwarded to the
// inner calls of an advisor, subagent or fusion server tool.
type ChatToolReasoning struct {
	Effort    ReasoningEffort `json:"effort,omitempty"`
	MaxTokens *int            `json:"max_tokens,omitempty"`
}

// ChatNestedTool represents a tool made available to the inner calls of an
// advisor, subagent or fusion server tool. Only OpenRouter server tools are
// supported, function tools are rejected.
type ChatNestedTool struct {
	Type       ChatToolType   `json:"type"`
	Parameters map[string]any `json:"parameters,omitempty"`
}

// ChatPlugin is implemented by every plugin configuration that can be enabled on
// a chat completion request.
type ChatPlugin interface {
	chatPlugin()
}

// ChatAutoRouterPlugin represents the auto-router plugin, routing between models
// by a cost and quality tradeoff.
type ChatAutoRouterPlugin struct {
	ID            ChatPluginID `json:"id"`
	AllowedModels []string     `json:"allowed_models,omitempty"`
	CostTier      ChatCostTier `json:"cost_tier,omitempty"`
	Enabled       *bool        `json:"enabled,omitempty"`
	PinModel      *bool        `json:"pin_model,omitempty"`

	// Deprecated: use CostTier instead. Takes precedence when both are set.
	CostQualityTradeoff *int `json:"cost_quality_tradeoff,omitempty"`
}

func (p ChatAutoRouterPlugin) chatPlugin() {}

// ChatAutoBetaRouterPlugin represents the auto-beta-router plugin, ranking models
// for the classified task type by community spend share.
type ChatAutoBetaRouterPlugin struct {
	ID            ChatPluginID `json:"id"`
	AllowedModels []string     `json:"allowed_models,omitempty"`
	CostTier      ChatCostTier `json:"cost_tier,omitempty"`
	Enabled       *bool        `json:"enabled,omitempty"`

	// Deprecated: use CostTier instead. Takes precedence when both are set.
	CostQualityTradeoff *int `json:"cost_quality_tradeoff,omitempty"`
}

func (p ChatAutoBetaRouterPlugin) chatPlugin() {}

// ChatContextCompressionPlugin represents the context-compression plugin.
type ChatContextCompressionPlugin struct {
	ID      ChatPluginID                 `json:"id"`
	Enabled *bool                        `json:"enabled,omitempty"`
	Engine  ChatContextCompressionEngine `json:"engine,omitempty"`
}

func (p ChatContextCompressionPlugin) chatPlugin() {}

// ChatFileParserPlugin represents the file-parser plugin.
type ChatFileParserPlugin struct {
	ID      ChatPluginID          `json:"id"`
	Enabled *bool                 `json:"enabled,omitempty"`
	PDF     *ChatPDFParserOptions `json:"pdf,omitempty"`
}

func (p ChatFileParserPlugin) chatPlugin() {}

// ChatPDFParserOptions holds the pdf parsing options of the file-parser plugin.
type ChatPDFParserOptions struct {
	Engine ChatPDFParserEngine `json:"engine,omitempty"`
}

// ChatFusionPlugin represents the fusion plugin, running an expert panel of
// models and synthesizing their answers. At most 8 analysis models are accepted,
// MaxToolCalls defaults to 8 and is capped at 16.
type ChatFusionPlugin struct {
	ID             ChatPluginID     `json:"id"`
	AnalysisModels []string         `json:"analysis_models,omitempty"`
	Enabled        *bool            `json:"enabled,omitempty"`
	MaxToolCalls   *int             `json:"max_tool_calls,omitempty"`
	Model          string           `json:"model,omitempty"`
	Preset         ChatFusionPreset `json:"preset,omitempty"`
	Tools          []ChatNestedTool `json:"tools,omitempty"`
}

func (p ChatFusionPlugin) chatPlugin() {}

// ChatModerationPlugin represents the moderation plugin.
type ChatModerationPlugin struct {
	ID ChatPluginID `json:"id"`
}

func (p ChatModerationPlugin) chatPlugin() {}

// ChatParetoRouterPlugin represents the pareto-router plugin. MaxPrice is in USD
// per million tokens and bypasses MinCodingScore when set.
type ChatParetoRouterPlugin struct {
	ID             ChatPluginID          `json:"id"`
	Enabled        *bool                 `json:"enabled,omitempty"`
	MaxPrice       *float64              `json:"max_price,omitempty"`
	MinCodingScore *float64              `json:"min_coding_score,omitempty"`
	PriceSource    ChatParetoPriceSource `json:"price_source,omitempty"`
}

func (p ChatParetoRouterPlugin) chatPlugin() {}

// ChatResponseHealingPlugin represents the response-healing plugin.
type ChatResponseHealingPlugin struct {
	ID      ChatPluginID `json:"id"`
	Enabled *bool        `json:"enabled,omitempty"`
}

func (p ChatResponseHealingPlugin) chatPlugin() {}

// ChatWebSearchPlugin represents the web search plugin.
type ChatWebSearchPlugin struct {
	ID             ChatPluginID               `json:"id"`
	Enabled        *bool                      `json:"enabled,omitempty"`
	Engine         ChatWebSearchEngine        `json:"engine,omitempty"`
	ExcludeDomains []string                   `json:"exclude_domains,omitempty"`
	IncludeDomains []string                   `json:"include_domains,omitempty"`
	MaxResults     *int                       `json:"max_results,omitempty"`
	MaxUses        *int                       `json:"max_uses,omitempty"`
	SearchPrompt   string                     `json:"search_prompt,omitempty"`
	UserLocation   *ChatWebSearchUserLocation `json:"user_location,omitempty"`
}

func (p ChatWebSearchPlugin) chatPlugin() {}

// ChatWebFetchPlugin represents the web-fetch plugin.
type ChatWebFetchPlugin struct {
	ID               ChatPluginID `json:"id"`
	AllowedDomains   []string     `json:"allowed_domains,omitempty"`
	BlockedDomains   []string     `json:"blocked_domains,omitempty"`
	MaxContentTokens *int         `json:"max_content_tokens,omitempty"`
	MaxUses          *int         `json:"max_uses,omitempty"`
}

func (p ChatWebFetchPlugin) chatPlugin() {}

// ChatCompletionResponse is the root response of the chat completion endpoint.
type ChatCompletionResponse struct {
	ID                 string              `json:"id"`
	Object             ChatObject          `json:"object"`
	Created            int64               `json:"created"`
	Model              string              `json:"model"`
	Choices            []ChatChoice        `json:"choices"`
	SystemFingerprint  *string             `json:"system_fingerprint"`
	ServiceTier        *string             `json:"service_tier"`
	Usage              *ChatUsage          `json:"usage,omitempty"`
	OpenRouterMetadata *OpenRouterMetadata `json:"openrouter_metadata,omitempty"`
}

// ChatChoice represents a single completion choice.
type ChatChoice struct {
	Index        int              `json:"index"`
	FinishReason ChatFinishReason `json:"finish_reason"`
	Message      ChatMessage      `json:"message"`
	Logprobs     *ChatLogprobs    `json:"logprobs,omitempty"`
}

// ChatStreamChunk represents a single chunk of a streaming chat completion. Error
// is populated when the request failed after the response started.
type ChatStreamChunk struct {
	ID                 string              `json:"id"`
	Object             ChatObject          `json:"object"`
	Created            int64               `json:"created"`
	Model              string              `json:"model"`
	Choices            []ChatStreamChoice  `json:"choices"`
	SystemFingerprint  string              `json:"system_fingerprint,omitempty"`
	ServiceTier        *string             `json:"service_tier,omitempty"`
	Usage              *ChatUsage          `json:"usage,omitempty"`
	OpenRouterMetadata *OpenRouterMetadata `json:"openrouter_metadata,omitempty"`
	Error              *ChatStreamError    `json:"error,omitempty"`
}

// ChatStreamChoice represents a single choice of a streaming chunk.
type ChatStreamChoice struct {
	Index        int              `json:"index"`
	FinishReason ChatFinishReason `json:"finish_reason"`
	Delta        ChatStreamDelta  `json:"delta"`
	Logprobs     *ChatLogprobs    `json:"logprobs,omitempty"`
}

// ChatStreamDelta represents the incremental changes of a streaming choice.
type ChatStreamDelta struct {
	Role             ChatRole              `json:"role,omitempty"`
	Content          string                `json:"content,omitempty"`
	Reasoning        string                `json:"reasoning,omitempty"`
	ReasoningDetails []ChatReasoningDetail `json:"reasoning_details,omitempty"`
	Refusal          string                `json:"refusal,omitempty"`
	Audio            *ChatAudioOutput      `json:"audio,omitempty"`
	Images           []ChatAssistantImage  `json:"images,omitempty"`
	ToolCalls        []ChatStreamToolCall  `json:"tool_calls,omitempty"`
}

// ChatStreamToolCall represents the delta of a tool call. Index identifies the
// tool call the delta belongs to, the remaining fields are only sent once the
// provider knows them.
type ChatStreamToolCall struct {
	Index    int                         `json:"index"`
	ID       string                      `json:"id,omitempty"`
	Type     ChatToolType                `json:"type,omitempty"`
	Function *ChatStreamToolCallFunction `json:"function,omitempty"`
}

// ChatStreamToolCallFunction holds the name and json encoded argument delta of a
// streamed tool call.
type ChatStreamToolCallFunction struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

// ChatStreamError represents the error details of a streaming completion that
// failed after the response started.
type ChatStreamError struct {
	Message  string                   `json:"message"`
	Code     int64                    `json:"code"`
	Metadata *ChatStreamErrorMetadata `json:"metadata,omitempty"`
}

// ChatStreamErrorMetadata represents the structured metadata of a streaming
// error.
type ChatStreamErrorMetadata struct {
	ErrorType    ChatErrorType `json:"error_type"`
	ProviderCode string        `json:"provider_code,omitempty"`
}

// ChatLogprobs represents the log probabilities of a completion.
type ChatLogprobs struct {
	Content []ChatTokenLogprob `json:"content"`
	Refusal []ChatTokenLogprob `json:"refusal"`
}

// ChatTokenLogprob represents the log probability of a single token.
type ChatTokenLogprob struct {
	Token       string           `json:"token"`
	Logprob     float64          `json:"logprob"`
	Bytes       []int            `json:"bytes"`
	TopLogprobs []ChatTopLogprob `json:"top_logprobs"`
}

// ChatTopLogprob represents a single alternative token and its log probability.
type ChatTopLogprob struct {
	Token   string  `json:"token"`
	Logprob float64 `json:"logprob"`
	Bytes   []int   `json:"bytes"`
}

// ChatUsage represents the token and cost usage of a chat completion.
type ChatUsage struct {
	PromptTokens            int                          `json:"prompt_tokens"`
	CompletionTokens        int                          `json:"completion_tokens"`
	TotalTokens             int                          `json:"total_tokens"`
	Cost                    *float64                     `json:"cost"`
	CostDetails             *CostDetails                 `json:"cost_details"`
	PromptTokensDetails     *PromptTokensDetails         `json:"prompt_tokens_details"`
	CompletionTokensDetails *ChatCompletionTokensDetails `json:"completion_tokens_details"`
	ServerToolUseDetails    *ServerToolUse               `json:"server_tool_use_details"`
	IsBYOK                  bool                         `json:"is_byok"`
}

// ChatCompletionTokensDetails represents the breakdown of tokens generated by
// the model, including the accepted and rejected tokens of a prediction.
type ChatCompletionTokensDetails struct {
	ReasoningTokens          *int `json:"reasoning_tokens"`
	AudioTokens              *int `json:"audio_tokens"`
	AcceptedPredictionTokens *int `json:"accepted_prediction_tokens"`
	RejectedPredictionTokens *int `json:"rejected_prediction_tokens"`
}

// OpenRouterMetadata represents the routing metadata of a request. It is only
// returned when ChatCompletionRequest.MetadataLevel is enabled.
type OpenRouterMetadata struct {
	Requested string            `json:"requested"`
	Strategy  RoutingStrategy   `json:"strategy"`
	Region    *string           `json:"region"`
	Summary   string            `json:"summary"`
	Attempt   int               `json:"attempt"`
	IsBYOK    bool              `json:"is_byok"`
	Endpoints EndpointsMetadata `json:"endpoints"`
	Attempts  []RouterAttempt   `json:"attempts,omitempty"`
	Params    *RouterParams     `json:"params,omitempty"`
	Pipeline  []PipelineStage   `json:"pipeline,omitempty"`
}

// EndpointsMetadata represents the endpoints considered while routing a request.
type EndpointsMetadata struct {
	Total     int            `json:"total"`
	Available []EndpointInfo `json:"available"`
}

// EndpointInfo represents a single endpoint considered while routing a request.
type EndpointInfo struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Selected bool   `json:"selected"`
}

// RouterAttempt represents a single upstream attempt of a request.
type RouterAttempt struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Status   int    `json:"status"`
}

// RouterParams represents the routing parameters a request was resolved with.
type RouterParams struct {
	QualityFloor    *float64 `json:"quality_floor,omitempty"`
	ThroughputFloor *float64 `json:"throughput_floor,omitempty"`
	VersionGroup    string   `json:"version_group,omitempty"`
}

// PipelineStage represents a single stage of the request pipeline. Multiple
// plugins share a type, Name disambiguates which one emitted the stage.
type PipelineStage struct {
	Type           PipelineStageType `json:"type"`
	Name           string            `json:"name"`
	Summary        string            `json:"summary,omitempty"`
	GuardrailID    string            `json:"guardrail_id,omitempty"`
	GuardrailScope string            `json:"guardrail_scope,omitempty"`
	CostUSD        *float64          `json:"cost_usd,omitempty"`
	Data           map[string]any    `json:"data,omitempty"`
}

// ChatMetadataLevel is the opt-in level for surfacing routing metadata on the
// response.
type ChatMetadataLevel string

const (
	ChatMetadataLevelDisabled ChatMetadataLevel = "disabled"
	ChatMetadataLevelEnabled  ChatMetadataLevel = "enabled"
)

// ChatRole is the role of the author of a message.
type ChatRole string

const (
	ChatRoleSystem    ChatRole = "system"
	ChatRoleDeveloper ChatRole = "developer"
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
	ChatRoleTool      ChatRole = "tool"
)

// ChatContentPartType is the type of a message content part.
type ChatContentPartType string

const (
	ChatContentPartTypeText       ChatContentPartType = "text"
	ChatContentPartTypeImageURL   ChatContentPartType = "image_url"
	ChatContentPartTypeInputAudio ChatContentPartType = "input_audio"
	ChatContentPartTypeVideoURL   ChatContentPartType = "video_url"
	ChatContentPartTypeFile       ChatContentPartType = "file"

	// Deprecated: use ChatContentPartTypeVideoURL instead.
	ChatContentPartTypeInputVideo ChatContentPartType = "input_video"
)

// ChatImageDetail is the detail level an image content part is processed with.
// Original is an OpenRouter extension requesting true original-resolution media,
// it is downgraded to high for providers without such a tier.
type ChatImageDetail string

const (
	ChatImageDetailAuto     ChatImageDetail = "auto"
	ChatImageDetailLow      ChatImageDetail = "low"
	ChatImageDetailHigh     ChatImageDetail = "high"
	ChatImageDetailOriginal ChatImageDetail = "original"
)

// AnthropicCacheControlType is the type of a cache control directive.
type AnthropicCacheControlType string

const (
	AnthropicCacheControlTypeEphemeral AnthropicCacheControlType = "ephemeral"
)

// AnthropicCacheTTL is the lifetime of a cache breakpoint.
type AnthropicCacheTTL string

const (
	AnthropicCacheTTL5M AnthropicCacheTTL = "5m"
	AnthropicCacheTTL1H AnthropicCacheTTL = "1h"
)

// ChatPromptCacheMode is the prompt caching mode of a request or content block.
type ChatPromptCacheMode string

const (
	ChatPromptCacheModeExplicit ChatPromptCacheMode = "explicit"
)

// ChatServiceTier is the service tier a request is processed with.
type ChatServiceTier string

const (
	ChatServiceTierAuto     ChatServiceTier = "auto"
	ChatServiceTierDefault  ChatServiceTier = "default"
	ChatServiceTierFlex     ChatServiceTier = "flex"
	ChatServiceTierPriority ChatServiceTier = "priority"
	ChatServiceTierScale    ChatServiceTier = "scale"
)

// ChatReasoningSummary is the verbosity of the reasoning summary of a request.
type ChatReasoningSummary string

const (
	ChatReasoningSummaryAuto     ChatReasoningSummary = "auto"
	ChatReasoningSummaryConcise  ChatReasoningSummary = "concise"
	ChatReasoningSummaryDetailed ChatReasoningSummary = "detailed"
)

// ChatPredictionType is the type of a predicted output.
type ChatPredictionType string

const (
	ChatPredictionTypeContent ChatPredictionType = "content"
)

// ChatResponseFormatType is the type of a response format configuration.
type ChatResponseFormatType string

const (
	ChatResponseFormatTypeText       ChatResponseFormatType = "text"
	ChatResponseFormatTypeJSONObject ChatResponseFormatType = "json_object"
	ChatResponseFormatTypeJSONSchema ChatResponseFormatType = "json_schema"
	ChatResponseFormatTypeGrammar    ChatResponseFormatType = "grammar"
	ChatResponseFormatTypePython     ChatResponseFormatType = "python"
)

// ChatRoute is the legacy alias of ProviderSortConfig.Partition.
//
// Deprecated: use ProviderPreferences.Sort.Partition instead.
type ChatRoute string

const (
	// Deprecated: use ProviderSortPartitionModel instead.
	ChatRouteFallback ChatRoute = "fallback"
	// Deprecated: use ProviderSortPartitionNone instead.
	ChatRouteSort ChatRoute = "sort"
)

// ChatStopConditionType is the type of a server tool stop condition.
type ChatStopConditionType string

const (
	ChatStopConditionTypeStepCountIs    ChatStopConditionType = "step_count_is"
	ChatStopConditionTypeHasToolCall    ChatStopConditionType = "has_tool_call"
	ChatStopConditionTypeMaxTokensUsed  ChatStopConditionType = "max_tokens_used"
	ChatStopConditionTypeMaxCost        ChatStopConditionType = "max_cost"
	ChatStopConditionTypeFinishReasonIs ChatStopConditionType = "finish_reason_is"
)

// ChatToolType is the type of a tool definition, tool call or forced tool choice.
type ChatToolType string

const (
	ChatToolTypeFunction        ChatToolType = "function"
	ChatToolTypeAdvisor         ChatToolType = "openrouter:advisor"
	ChatToolTypeBash            ChatToolType = "openrouter:bash"
	ChatToolTypeDatetime        ChatToolType = "openrouter:datetime"
	ChatToolTypeFiles           ChatToolType = "openrouter:files"
	ChatToolTypeFusion          ChatToolType = "openrouter:fusion"
	ChatToolTypeImageGeneration ChatToolType = "openrouter:image_generation"
	ChatToolTypeSearchModels    ChatToolType = "openrouter:experimental__search_models"
	ChatToolTypeSubagent        ChatToolType = "openrouter:subagent"
	ChatToolTypeWebFetch        ChatToolType = "openrouter:web_fetch"
	ChatToolTypeWebSearch       ChatToolType = "openrouter:web_search"

	ChatToolTypeWebSearchShorthand       ChatToolType = "web_search"
	ChatToolTypeWebSearchPreview         ChatToolType = "web_search_preview"
	ChatToolTypeWebSearchPreview20250311 ChatToolType = "web_search_preview_2025_03_11"
	ChatToolTypeWebSearch20250826        ChatToolType = "web_search_2025_08_26"
)

// ChatToolChoiceMode is a plain tool choice, without naming a specific tool.
type ChatToolChoiceMode string

const (
	ChatToolChoiceModeNone     ChatToolChoiceMode = "none"
	ChatToolChoiceModeAuto     ChatToolChoiceMode = "auto"
	ChatToolChoiceModeRequired ChatToolChoiceMode = "required"
)

// ChatBashEngine is the engine backing the bash server tool. Auto and native
// return the tool call to the caller, openrouter executes it server-side.
type ChatBashEngine string

const (
	ChatBashEngineAuto       ChatBashEngine = "auto"
	ChatBashEngineNative     ChatBashEngine = "native"
	ChatBashEngineOpenRouter ChatBashEngine = "openrouter"
)

// ChatBashEnvironmentType is the kind of execution environment of the bash
// server tool.
type ChatBashEnvironmentType string

const (
	ChatBashEnvironmentTypeContainerAuto      ChatBashEnvironmentType = "container_auto"
	ChatBashEnvironmentTypeContainerReference ChatBashEnvironmentType = "container_reference"
)

// ChatWebSearchEngine is the engine backing web search. Auto is only valid on
// server tool configurations, the web search plugin requires an explicit engine.
type ChatWebSearchEngine string

const (
	ChatWebSearchEngineAuto       ChatWebSearchEngine = "auto"
	ChatWebSearchEngineNative     ChatWebSearchEngine = "native"
	ChatWebSearchEngineExa        ChatWebSearchEngine = "exa"
	ChatWebSearchEngineParallel   ChatWebSearchEngine = "parallel"
	ChatWebSearchEngineFirecrawl  ChatWebSearchEngine = "firecrawl"
	ChatWebSearchEnginePerplexity ChatWebSearchEngine = "perplexity"
)

// ChatWebFetchEngine is the engine backing web fetch.
type ChatWebFetchEngine string

const (
	ChatWebFetchEngineAuto       ChatWebFetchEngine = "auto"
	ChatWebFetchEngineNative     ChatWebFetchEngine = "native"
	ChatWebFetchEngineOpenRouter ChatWebFetchEngine = "openrouter"
	ChatWebFetchEngineExa        ChatWebFetchEngine = "exa"
	ChatWebFetchEngineParallel   ChatWebFetchEngine = "parallel"
	ChatWebFetchEngineFirecrawl  ChatWebFetchEngine = "firecrawl"
)

// ChatSearchContextSize is how much context is retrieved per search result. It
// is overridden by an explicit character cap.
type ChatSearchContextSize string

const (
	ChatSearchContextSizeLow    ChatSearchContextSize = "low"
	ChatSearchContextSizeMedium ChatSearchContextSize = "medium"
	ChatSearchContextSizeHigh   ChatSearchContextSize = "high"
)

// ChatUserLocationType is the precision of a web search user location.
type ChatUserLocationType string

const (
	ChatUserLocationTypeApproximate ChatUserLocationType = "approximate"
)

// ChatPluginID is the identifier of a plugin.
type ChatPluginID string

const (
	ChatPluginIDAutoRouter         ChatPluginID = "auto-router"
	ChatPluginIDAutoBetaRouter     ChatPluginID = "auto-beta-router"
	ChatPluginIDModeration         ChatPluginID = "moderation"
	ChatPluginIDWeb                ChatPluginID = "web"
	ChatPluginIDWebFetch           ChatPluginID = "web-fetch"
	ChatPluginIDFileParser         ChatPluginID = "file-parser"
	ChatPluginIDResponseHealing    ChatPluginID = "response-healing"
	ChatPluginIDContextCompression ChatPluginID = "context-compression"
	ChatPluginIDParetoRouter       ChatPluginID = "pareto-router"
	ChatPluginIDFusion             ChatPluginID = "fusion"
)

// ChatCostTier is the named cost and quality setting of a router plugin.
type ChatCostTier string

const (
	ChatCostTierLow    ChatCostTier = "low"
	ChatCostTierMedium ChatCostTier = "medium"
	ChatCostTierHigh   ChatCostTier = "high"
	ChatCostTierXHigh  ChatCostTier = "xhigh"
	ChatCostTierMax    ChatCostTier = "max"
)

// ChatContextCompressionEngine is the engine used to compress the context.
type ChatContextCompressionEngine string

const (
	ChatContextCompressionEngineMiddleOut ChatContextCompressionEngine = "middle-out"
)

// ChatPDFParserEngine is the engine used to parse pdf files.
type ChatPDFParserEngine string

const (
	ChatPDFParserEngineMistralOCR   ChatPDFParserEngine = "mistral-ocr"
	ChatPDFParserEngineNative       ChatPDFParserEngine = "native"
	ChatPDFParserEngineCloudflareAI ChatPDFParserEngine = "cloudflare-ai"

	// Deprecated: automatically redirected to ChatPDFParserEngineCloudflareAI.
	ChatPDFParserEnginePDFText ChatPDFParserEngine = "pdf-text"
)

// ChatFusionPreset is a curated panel and analyst configuration of the fusion
// plugin.
type ChatFusionPreset string

const (
	ChatFusionPresetGeneralHigh   ChatFusionPreset = "general-high"
	ChatFusionPresetGeneralBudget ChatFusionPreset = "general-budget"
	ChatFusionPresetGeneralFast   ChatFusionPreset = "general-fast"
)

// ChatParetoPriceSource is the price used as the cost axis of the pareto router.
type ChatParetoPriceSource string

const (
	ChatParetoPriceSourcePrompt      ChatParetoPriceSource = "prompt"
	ChatParetoPriceSourceWeightedAvg ChatParetoPriceSource = "weighted_avg"
)

// ChatObject is the object type of a chat completion response or chunk.
type ChatObject string

const (
	ChatObjectCompletion      ChatObject = "chat.completion"
	ChatObjectCompletionChunk ChatObject = "chat.completion.chunk"
)

// ChatFinishReason is the reason a completion stopped.
type ChatFinishReason string

const (
	ChatFinishReasonToolCalls     ChatFinishReason = "tool_calls"
	ChatFinishReasonStop          ChatFinishReason = "stop"
	ChatFinishReasonLength        ChatFinishReason = "length"
	ChatFinishReasonContentFilter ChatFinishReason = "content_filter"
	ChatFinishReasonError         ChatFinishReason = "error"
)

// ChatReasoningDetailType is the type of a reasoning detail.
type ChatReasoningDetailType string

const (
	ChatReasoningDetailTypeSummary        ChatReasoningDetailType = "reasoning.summary"
	ChatReasoningDetailTypeEncrypted      ChatReasoningDetailType = "reasoning.encrypted"
	ChatReasoningDetailTypeText           ChatReasoningDetailType = "reasoning.text"
	ChatReasoningDetailTypeServerToolCall ChatReasoningDetailType = "reasoning.server_tool_call"
)

// ChatReasoningFormat is the upstream format a reasoning detail was produced in.
type ChatReasoningFormat string

const (
	ChatReasoningFormatUnknown                  ChatReasoningFormat = "unknown"
	ChatReasoningFormatOpenAIResponsesV1        ChatReasoningFormat = "openai-responses-v1"
	ChatReasoningFormatAzureOpenAIResponsesV1   ChatReasoningFormat = "azure-openai-responses-v1"
	ChatReasoningFormatBedrockOpenAIResponsesV1 ChatReasoningFormat = "bedrock-openai-responses-v1"
	ChatReasoningFormatXAIResponsesV1           ChatReasoningFormat = "xai-responses-v1"
	ChatReasoningFormatMetaResponsesV1          ChatReasoningFormat = "meta-responses-v1"
	ChatReasoningFormatAnthropicClaudeV1        ChatReasoningFormat = "anthropic-claude-v1"
	ChatReasoningFormatGoogleGeminiV1           ChatReasoningFormat = "google-gemini-v1"
)

// ChatErrorType is the canonical OpenRouter error type of a streaming error, it
// is stable across all api formats.
type ChatErrorType string

const (
	ChatErrorTypeContextLengthExceeded  ChatErrorType = "context_length_exceeded"
	ChatErrorTypeMaxTokensExceeded      ChatErrorType = "max_tokens_exceeded"
	ChatErrorTypeTokenLimitExceeded     ChatErrorType = "token_limit_exceeded"
	ChatErrorTypeStringTooLong          ChatErrorType = "string_too_long"
	ChatErrorTypeAuthentication         ChatErrorType = "authentication"
	ChatErrorTypePermissionDenied       ChatErrorType = "permission_denied"
	ChatErrorTypePaymentRequired        ChatErrorType = "payment_required"
	ChatErrorTypeRateLimitExceeded      ChatErrorType = "rate_limit_exceeded"
	ChatErrorTypeProviderOverloaded     ChatErrorType = "provider_overloaded"
	ChatErrorTypeProviderUnavailable    ChatErrorType = "provider_unavailable"
	ChatErrorTypeInvalidRequest         ChatErrorType = "invalid_request"
	ChatErrorTypeInvalidPrompt          ChatErrorType = "invalid_prompt"
	ChatErrorTypeNotFound               ChatErrorType = "not_found"
	ChatErrorTypePreconditionFailed     ChatErrorType = "precondition_failed"
	ChatErrorTypePayloadTooLarge        ChatErrorType = "payload_too_large"
	ChatErrorTypeUnprocessable          ChatErrorType = "unprocessable"
	ChatErrorTypeContentPolicyViolation ChatErrorType = "content_policy_violation"
	ChatErrorTypeRefusal                ChatErrorType = "refusal"
	ChatErrorTypeInvalidImage           ChatErrorType = "invalid_image"
	ChatErrorTypeImageTooLarge          ChatErrorType = "image_too_large"
	ChatErrorTypeImageTooSmall          ChatErrorType = "image_too_small"
	ChatErrorTypeUnsupportedImageFormat ChatErrorType = "unsupported_image_format"
	ChatErrorTypeImageNotFound          ChatErrorType = "image_not_found"
	ChatErrorTypeImageDownloadFailed    ChatErrorType = "image_download_failed"
	ChatErrorTypeServer                 ChatErrorType = "server"
	ChatErrorTypeTimeout                ChatErrorType = "timeout"
	ChatErrorTypeUnmapped               ChatErrorType = "unmapped"
)

// RoutingStrategy is the strategy a request was routed with.
type RoutingStrategy string

const (
	RoutingStrategyDirect      RoutingStrategy = "direct"
	RoutingStrategyAuto        RoutingStrategy = "auto"
	RoutingStrategyFree        RoutingStrategy = "free"
	RoutingStrategyLatest      RoutingStrategy = "latest"
	RoutingStrategyAlias       RoutingStrategy = "alias"
	RoutingStrategyFallback    RoutingStrategy = "fallback"
	RoutingStrategyPareto      RoutingStrategy = "pareto"
	RoutingStrategyBodybuilder RoutingStrategy = "bodybuilder"
	RoutingStrategyFusion      RoutingStrategy = "fusion"
)

// PipelineStageType is the categorical kind of a pipeline stage.
type PipelineStageType string

const (
	PipelineStageTypeGuardrail          PipelineStageType = "guardrail"
	PipelineStageTypePlugin             PipelineStageType = "plugin"
	PipelineStageTypeServerTools        PipelineStageType = "server_tools"
	PipelineStageTypeResponseHealing    PipelineStageType = "response_healing"
	PipelineStageTypeContextCompression PipelineStageType = "context_compression"
)
