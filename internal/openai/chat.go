package openai

// ChatObject is the object type of a chat completion response or chunk.
type ChatObject string

const (
	ChatObjectCompletion      ChatObject = "chat.completion"
	ChatObjectCompletionChunk ChatObject = "chat.completion.chunk"
)

// ChatRole is the role of the author of a message.
type ChatRole string

const (
	ChatRoleDeveloper ChatRole = "developer"
	ChatRoleSystem    ChatRole = "system"
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
	ChatRoleTool      ChatRole = "tool"
	ChatRoleFunction  ChatRole = "function"
)

// ChatContentPartType is the type of a message content part.
type ChatContentPartType string

const (
	ChatContentPartTypeText       ChatContentPartType = "text"
	ChatContentPartTypeImageURL   ChatContentPartType = "image_url"
	ChatContentPartTypeInputAudio ChatContentPartType = "input_audio"
	ChatContentPartTypeFile       ChatContentPartType = "file"
)

// ChatCompletionRequest is the request body of the chat completions endpoint.
// Messages is required, Model may be omitted when the server decides the route.
// Zero values and nil pointers of the remaining fields are omitted.
type ChatCompletionRequest struct {
	Model    string      `json:"model,omitempty"`
	Messages []ChatMessage `json:"messages"`

	FrequencyPenalty    *float64                `json:"frequency_penalty,omitempty"`
	LogitBias           map[string]float64      `json:"logit_bias,omitempty"`
	Logprobs            *bool                   `json:"logprobs,omitempty"`
	MaxCompletionTokens *int                    `json:"max_completion_tokens,omitempty"`
	MaxTokens           *int                    `json:"max_tokens,omitempty"`
	Metadata            map[string]string       `json:"metadata,omitempty"`
	Modalities          []string                `json:"modalities,omitempty"`
	ParallelToolCalls   *bool                   `json:"parallel_tool_calls,omitempty"`
	Prediction          *ChatPrediction         `json:"prediction,omitempty"`
	PresencePenalty     *float64                `json:"presence_penalty,omitempty"`
	PromptCacheKey      string                  `json:"prompt_cache_key,omitempty"`
	PromptCacheOptions  *ChatPromptCacheOptions `json:"prompt_cache_options,omitempty"`
	ReasoningEffort     string                  `json:"reasoning_effort,omitempty"`
	ResponseFormat      *ChatResponseFormat     `json:"response_format,omitempty"`
	Seed                *int                    `json:"seed,omitempty"`
	ServiceTier         string                  `json:"service_tier,omitempty"`
	Stop                []string                `json:"stop,omitempty"`
	Stream              *bool                   `json:"stream,omitempty"`
	StreamOptions       *ChatStreamOptions      `json:"stream_options,omitempty"`
	Temperature         *float64                `json:"temperature,omitempty"`
	ToolChoice          *ChatToolChoice         `json:"tool_choice,omitempty"`
	Tools               []ChatFunctionTool      `json:"tools,omitempty"`
	TopLogprobs         *int                    `json:"top_logprobs,omitempty"`
	TopP                *float64                `json:"top_p,omitempty"`
	User                string                  `json:"user,omitempty"`
	Verbosity           string                  `json:"verbosity,omitempty"`
}

// ChatMessage represents a single message of a conversation. Role determines
// which of the remaining fields are meaningful. Content is a string or a list
// of content parts.
type ChatMessage struct {
	Role    ChatRole    `json:"role"`
	Content *ChatContent `json:"content,omitempty"`
	Name    string      `json:"name,omitempty"`

	Refusal    string       `json:"refusal,omitempty"`
	Audio      *ChatAudio   `json:"audio,omitempty"`
	ToolCalls  []ChatToolCall `json:"tool_calls,omitempty"`
	ToolCallID string       `json:"tool_call_id,omitempty"`
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
		return jsonMarshal(cc.Parts)
	}

	return jsonMarshal(cc.Text)
}

// ChatContentPart represents a single content part of a message. Type
// determines which of the remaining fields are used.
type ChatContentPart struct {
	Type       ChatContentPartType  `json:"type"`
	Text       string               `json:"text,omitempty"`
	ImageURL   *ChatContentImageURL `json:"image_url,omitempty"`
	InputAudio *ChatContentInputAudio `json:"input_audio,omitempty"`
	File       *ChatContentFile     `json:"file,omitempty"`
}

// ChatContentImageURL holds the url of an image content part, as a base64 data
// url or an HTTP(S) url.
type ChatContentImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

// ChatContentInputAudio holds the base64 encoded audio of an audio content part.
type ChatContentInputAudio struct {
	Data   string `json:"data"`
	Format string `json:"format"`
}

// ChatContentFile holds the document of a file content part, either inline as a
// base64 data url or url, or by the id of a previously uploaded file.
type ChatContentFile struct {
	FileData string `json:"file_data,omitempty"`
	FileID   string `json:"file_id,omitempty"`
	Filename string `json:"filename,omitempty"`
}

// ChatAudio references an audio output of a previous assistant message.
type ChatAudio struct {
	ID string `json:"id,omitempty"`
}

// ChatToolCall represents a tool call made by the assistant.
type ChatToolCall struct {
	ID       string               `json:"id"`
	Type     string               `json:"type"`
	Function ChatToolCallFunction `json:"function"`
}

// ChatToolCallFunction holds the name and json encoded arguments of a tool call.
type ChatToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ChatFunctionTool represents a regular function tool.
type ChatFunctionTool struct {
	Type     string       `json:"type"`
	Function ChatFunction `json:"function"`
}

// ChatFunction represents the definition of a callable function.
type ChatFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
	Strict      *bool          `json:"strict,omitempty"`
}

// ChatToolChoice represents the tool choice configuration of a request. Mode
// encodes the plain "none", "auto" and "required" choices and takes precedence,
// otherwise Function names the tool to force.
type ChatToolChoice struct {
	Mode     string                `json:"-"`
	Type     string                `json:"type,omitempty"`
	Function *ChatToolChoiceFunc   `json:"function,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface for ChatToolChoice.
func (tc ChatToolChoice) MarshalJSON() ([]byte, error) {
	if tc.Mode != "" {
		return jsonMarshal(tc.Mode)
	}

	type choice ChatToolChoice

	return jsonMarshal(choice(tc))
}

// ChatToolChoiceFunc holds the name of the function a named tool choice forces.
type ChatToolChoiceFunc struct {
	Name string `json:"name"`
}

// ChatResponseFormat represents the response format configuration of a request.
// Type determines which of the remaining fields are used.
type ChatResponseFormat struct {
	Type       string         `json:"type"`
	JSONSchema *ChatJSONSchema `json:"json_schema,omitempty"`
}

// ChatJSONSchema represents the schema of a json_schema response format. Name is
// required.
type ChatJSONSchema struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Schema      map[string]any `json:"schema,omitempty"`
	Strict      *bool          `json:"strict,omitempty"`
}

// ChatPrediction represents static predicted output content.
type ChatPrediction struct {
	Type    string       `json:"type"`
	Content ChatContent  `json:"content"`
}

// ChatPromptCacheOptions represents the request level prompt cache controls.
type ChatPromptCacheOptions struct {
	Mode string `json:"mode,omitempty"`
	TTL  string `json:"ttl,omitempty"`
}

// ChatStreamOptions holds the streaming options of a request.
type ChatStreamOptions struct {
	IncludeUsage *bool `json:"include_usage,omitempty"`
}

// ChatCompletionResponse is the root response of the chat completions endpoint.
type ChatCompletionResponse struct {
	ID                string         `json:"id"`
	Object            ChatObject     `json:"object"`
	Created           int64          `json:"created"`
	Model             string         `json:"model"`
	Choices           []ChatChoice   `json:"choices"`
	SystemFingerprint *string        `json:"system_fingerprint"`
	ServiceTier       string         `json:"service_tier,omitempty"`
	Usage             *ChatUsage     `json:"usage,omitempty"`
}

// ChatChoice represents a single completion choice.
type ChatChoice struct {
	Index        int           `json:"index"`
	FinishReason string        `json:"finish_reason"`
	Message      ChatOutMessage `json:"message"`
	Logprobs     *ChatLogprobs `json:"logprobs,omitempty"`
}

// ChatOutMessage represents the assistant message of a completion choice.
type ChatOutMessage struct {
	Role      string        `json:"role"`
	Content   *string       `json:"content"`
	Refusal   string        `json:"refusal,omitempty"`
	Audio     *ChatOutAudio `json:"audio,omitempty"`
	ToolCalls []ChatToolCall `json:"tool_calls,omitempty"`
}

// ChatOutAudio represents the audio output of an assistant message.
type ChatOutAudio struct {
	ID         string `json:"id,omitempty"`
	Data       string `json:"data,omitempty"`
	Transcript string `json:"transcript,omitempty"`
	ExpiresAt  int64  `json:"expires_at,omitempty"`
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

// ChatUsage represents the token usage of a chat completion.
type ChatUsage struct {
	PromptTokens            int                        `json:"prompt_tokens"`
	CompletionTokens        int                        `json:"completion_tokens"`
	TotalTokens             int                        `json:"total_tokens"`
	PromptTokensDetails     *ChatPromptTokensDetails   `json:"prompt_tokens_details,omitempty"`
	CompletionTokensDetails *ChatCompletionTokensDetails `json:"completion_tokens_details,omitempty"`
}

// ChatPromptTokensDetails represents the breakdown of tokens used in the prompt.
type ChatPromptTokensDetails struct {
	CachedTokens     *int `json:"cached_tokens,omitempty"`
	CacheWriteTokens *int `json:"cache_write_tokens,omitempty"`
	AudioTokens      *int `json:"audio_tokens,omitempty"`
	ImageTokens      *int `json:"image_tokens,omitempty"`
	TextTokens       *int `json:"text_tokens,omitempty"`
}

// ChatCompletionTokensDetails represents the breakdown of tokens generated by
// the model.
type ChatCompletionTokensDetails struct {
	ReasoningTokens          *int `json:"reasoning_tokens,omitempty"`
	AudioTokens              *int `json:"audio_tokens,omitempty"`
	AcceptedPredictionTokens *int `json:"accepted_prediction_tokens,omitempty"`
	RejectedPredictionTokens *int `json:"rejected_prediction_tokens,omitempty"`
}

// ChatCompletionChunk represents a single chunk of a streaming chat completion.
type ChatCompletionChunk struct {
	ID                string            `json:"id"`
	Object            ChatObject        `json:"object"`
	Created           int64             `json:"created"`
	Model             string            `json:"model"`
	Choices           []ChatChunkChoice `json:"choices"`
	SystemFingerprint *string           `json:"system_fingerprint"`
	ServiceTier       string            `json:"service_tier,omitempty"`
	Usage             *ChatUsage        `json:"usage,omitempty"`
}

// ChatChunkChoice represents a single choice of a streaming chunk.
type ChatChunkChoice struct {
	Index        int           `json:"index"`
	FinishReason *string       `json:"finish_reason"`
	Delta        ChatChunkDelta `json:"delta"`
	Logprobs     *ChatLogprobs `json:"logprobs,omitempty"`
}

// ChatChunkDelta represents the incremental changes of a streaming choice.
type ChatChunkDelta struct {
	Role      string             `json:"role,omitempty"`
	Content   string             `json:"content,omitempty"`
	Refusal   string             `json:"refusal,omitempty"`
	Audio     *ChatOutAudio      `json:"audio,omitempty"`
	ToolCalls []ChatStreamToolCall `json:"tool_calls,omitempty"`
}

// ChatStreamToolCall represents the delta of a tool call. Index identifies the
// tool call the delta belongs to, the remaining fields are only sent once the
// provider knows them.
type ChatStreamToolCall struct {
	Index    int                         `json:"index"`
	ID       string                      `json:"id,omitempty"`
	Type     string                      `json:"type,omitempty"`
	Function *ChatStreamToolCallFunction `json:"function,omitempty"`
}

// ChatStreamToolCallFunction holds the name and json encoded argument delta of a
// streamed tool call.
type ChatStreamToolCallFunction struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}
