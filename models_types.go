package openingrouter

// ModelsListResponse is the root response for the models list endpoints.
type ModelsListResponse struct {
	Data       []Model         `json:"data"`
	TotalCount int             `json:"total_count"`
	Links      ModelsListLinks `json:"links"`
}

// ModelsListLinks represents the pagination links of a models list response.
type ModelsListLinks struct {
	Next *string `json:"next"`
}

// ModelResponse is the root response for the single model endpoint.
type ModelResponse struct {
	Data Model `json:"data"`
}

// Model represents an AI model available on OpenRouter.
type Model struct {
	ID                  string                  `json:"id"`
	CanonicalSlug       string                  `json:"canonical_slug"`
	Name                string                  `json:"name"`
	Description         string                  `json:"description"`
	Created             int64                   `json:"created"`
	ContextLength       *int                    `json:"context_length"`
	HuggingFaceID       *string                 `json:"hugging_face_id,omitempty"`
	ExpirationDate      *FlexibleTime           `json:"expiration_date"`
	KnowledgeCutoff     *FlexibleTime           `json:"knowledge_cutoff"`
	Architecture        ModelArchitecture       `json:"architecture"`
	Pricing             ModelPricing            `json:"pricing"`
	TopProvider         ModelTopProvider        `json:"top_provider"`
	PerRequestLimits    *ModelPerRequestLimits  `json:"per_request_limits"`
	DefaultParameters   *ModelDefaultParameters `json:"default_parameters"`
	SupportedParameters []Parameter             `json:"supported_parameters"`
	SupportedVoices     []string                `json:"supported_voices"`
	Reasoning           *ModelReasoning         `json:"reasoning,omitempty"`
	Benchmarks          *ModelBenchmarks        `json:"benchmarks,omitempty"`
	Links               ModelLinks              `json:"links"`
}

// ModelArchitecture represents the architecture information of a model.
type ModelArchitecture struct {
	Modality         *string          `json:"modality"`
	InputModalities  []InputModality  `json:"input_modalities"`
	OutputModalities []OutputModality `json:"output_modalities"`
	InstructType     *InstructType    `json:"instruct_type"`
	Tokenizer        ModelGroup       `json:"tokenizer,omitempty"`
}

// ModelLinks represents related api endpoints and resources of a model.
type ModelLinks struct {
	Details string `json:"details"`
}

// ModelTopProvider represents information about the top provider of a model.
type ModelTopProvider struct {
	ContextLength       *int `json:"context_length"`
	MaxCompletionTokens *int `json:"max_completion_tokens"`
	IsModerated         bool `json:"is_moderated"`
}

// ModelPerRequestLimits represents the per-request token limits of a model.
type ModelPerRequestLimits struct {
	PromptTokens     StringifiedNumber `json:"prompt_tokens"`
	CompletionTokens StringifiedNumber `json:"completion_tokens"`
}

// ModelDefaultParameters represents the default sampling parameters of a model.
type ModelDefaultParameters struct {
	Temperature       *float64 `json:"temperature"`
	TopP              *float64 `json:"top_p"`
	TopK              *int     `json:"top_k"`
	FrequencyPenalty  *float64 `json:"frequency_penalty"`
	PresencePenalty   *float64 `json:"presence_penalty"`
	RepetitionPenalty *float64 `json:"repetition_penalty"`
}

// ModelPricing represents the pricing information of a model. All prices are in
// USD, per token unless documented otherwise.
type ModelPricing struct {
	Prompt            StringifiedNumber      `json:"prompt"`
	Completion        StringifiedNumber      `json:"completion"`
	Request           StringifiedNumber      `json:"request,omitempty"`
	Image             StringifiedNumber      `json:"image,omitempty"`
	ImageOutput       StringifiedNumber      `json:"image_output,omitempty"`
	ImageToken        StringifiedNumber      `json:"image_token,omitempty"`
	Audio             StringifiedNumber      `json:"audio,omitempty"`
	AudioOutput       StringifiedNumber      `json:"audio_output,omitempty"`
	InputAudioCache   StringifiedNumber      `json:"input_audio_cache,omitempty"`
	InputCacheRead    StringifiedNumber      `json:"input_cache_read,omitempty"`
	InputCacheWrite   StringifiedNumber      `json:"input_cache_write,omitempty"`
	InputCacheWrite1H StringifiedNumber      `json:"input_cache_write_1h,omitempty"`
	InternalReasoning StringifiedNumber      `json:"internal_reasoning,omitempty"`
	WebSearch         StringifiedNumber      `json:"web_search,omitempty"`
	Discount          float64                `json:"discount"`
	Overrides         []ModelPricingOverride `json:"overrides,omitempty"`
}

// ModelPricingOverride represents a conditional override of the base pricing. An
// entry applies only when all of its condition fields match the request, among
// applicable entries later entries win per price key and keys absent from an
// entry inherit the base price.
type ModelPricingOverride struct {
	Prompt            StringifiedNumber `json:"prompt,omitempty"`
	Completion        StringifiedNumber `json:"completion,omitempty"`
	Audio             StringifiedNumber `json:"audio,omitempty"`
	InputAudioCache   StringifiedNumber `json:"input_audio_cache,omitempty"`
	InputCacheRead    StringifiedNumber `json:"input_cache_read,omitempty"`
	InputCacheWrite   StringifiedNumber `json:"input_cache_write,omitempty"`
	InputCacheWrite1H StringifiedNumber `json:"input_cache_write_1h,omitempty"`
	MinPromptTokens   *float64          `json:"min_prompt_tokens,omitempty"`
	UTCStart          *float64          `json:"utc_start,omitempty"`
	UTCEnd            *float64          `json:"utc_end,omitempty"`
}

// ModelReasoning represents the reasoning effort configuration of a model. It is
// omitted for non-reasoning models and dynamic router models.
type ModelReasoning struct {
	Mandatory         bool              `json:"mandatory"`
	DefaultEffort     *ReasoningEffort  `json:"default_effort,omitempty"`
	DefaultEnabled    *bool             `json:"default_enabled,omitempty"`
	SupportedEfforts  []ReasoningEffort `json:"supported_efforts"`
	SupportsMaxTokens *bool             `json:"supports_max_tokens,omitempty"`
}

// ModelBenchmarks represents third-party benchmark rankings of a model. It is
// omitted when no benchmark data is available.
type ModelBenchmarks struct {
	ArtificialAnalysis *ArtificialAnalysisBenchmark `json:"artificial_analysis,omitempty"`
	DesignArena        []DesignArenaBenchmark       `json:"design_arena"`
}

// ArtificialAnalysisBenchmark represents the Artificial Analysis index scores.
type ArtificialAnalysisBenchmark struct {
	IntelligenceIndex *float64 `json:"intelligence_index"`
	CodingIndex       *float64 `json:"coding_index"`
	AgenticIndex      *float64 `json:"agentic_index"`
}

// DesignArenaBenchmark represents a single Design Arena entry for a specific
// arena and category pair.
type DesignArenaBenchmark struct {
	Arena    string  `json:"arena"`
	Category string  `json:"category"`
	ELO      float64 `json:"elo"`
	WinRate  float64 `json:"win_rate"`
	Rank     int     `json:"rank"`
}

// Parameter is a request parameter that may be supported by a model.
type Parameter string

const (
	ParameterTemperature         Parameter = "temperature"
	ParameterTopP                Parameter = "top_p"
	ParameterTopK                Parameter = "top_k"
	ParameterMinP                Parameter = "min_p"
	ParameterTopA                Parameter = "top_a"
	ParameterFrequencyPenalty    Parameter = "frequency_penalty"
	ParameterPresencePenalty     Parameter = "presence_penalty"
	ParameterRepetitionPenalty   Parameter = "repetition_penalty"
	ParameterMaxTokens           Parameter = "max_tokens"
	ParameterMaxCompletionTokens Parameter = "max_completion_tokens"
	ParameterLogitBias           Parameter = "logit_bias"
	ParameterLogprobs            Parameter = "logprobs"
	ParameterTopLogprobs         Parameter = "top_logprobs"
	ParameterPrediction          Parameter = "prediction"
	ParameterSeed                Parameter = "seed"
	ParameterResponseFormat      Parameter = "response_format"
	ParameterStructuredOutputs   Parameter = "structured_outputs"
	ParameterStop                Parameter = "stop"
	ParameterTools               Parameter = "tools"
	ParameterToolChoice          Parameter = "tool_choice"
	ParameterParallelToolCalls   Parameter = "parallel_tool_calls"
	ParameterIncludeReasoning    Parameter = "include_reasoning"
	ParameterReasoning           Parameter = "reasoning"
	ParameterReasoningEffort     Parameter = "reasoning_effort"
	ParameterWebSearchOptions    Parameter = "web_search_options"
	ParameterVerbosity           Parameter = "verbosity"
)

// InputModality is a modality a model accepts as input.
type InputModality string

const (
	InputModalityText  InputModality = "text"
	InputModalityImage InputModality = "image"
	InputModalityFile  InputModality = "file"
	InputModalityAudio InputModality = "audio"
	InputModalityVideo InputModality = "video"
)

// OutputModality is a modality a model produces as output.
type OutputModality string

const (
	OutputModalityText          OutputModality = "text"
	OutputModalityImage         OutputModality = "image"
	OutputModalityEmbeddings    OutputModality = "embeddings"
	OutputModalityAudio         OutputModality = "audio"
	OutputModalityVideo         OutputModality = "video"
	OutputModalityRerank        OutputModality = "rerank"
	OutputModalitySpeech        OutputModality = "speech"
	OutputModalityTranscription OutputModality = "transcription"
)

// InstructType is the instruction format type of a model.
type InstructType string

const (
	InstructTypeNone        InstructType = "none"
	InstructTypeAiroboros   InstructType = "airoboros"
	InstructTypeAlpaca      InstructType = "alpaca"
	InstructTypeAlpacaModif InstructType = "alpaca-modif"
	InstructTypeChatML      InstructType = "chatml"
	InstructTypeClaude      InstructType = "claude"
	InstructTypeCodeLlama   InstructType = "code-llama"
	InstructTypeGemma       InstructType = "gemma"
	InstructTypeLlama2      InstructType = "llama2"
	InstructTypeLlama3      InstructType = "llama3"
	InstructTypeMistral     InstructType = "mistral"
	InstructTypeNemotron    InstructType = "nemotron"
	InstructTypeNeural      InstructType = "neural"
	InstructTypeOpenChat    InstructType = "openchat"
	InstructTypePhi3        InstructType = "phi3"
	InstructTypeRWKV        InstructType = "rwkv"
	InstructTypeVicuna      InstructType = "vicuna"
	InstructTypeZephyr      InstructType = "zephyr"
	InstructTypeDeepSeekR1  InstructType = "deepseek-r1"
	InstructTypeDeepSeekV31 InstructType = "deepseek-v3.1"
	InstructTypeQwQ         InstructType = "qwq"
	InstructTypeQwen3       InstructType = "qwen3"
)

// ModelGroup is the tokenizer type (model family) used by a model.
type ModelGroup string

const (
	ModelGroupRouter   ModelGroup = "Router"
	ModelGroupMedia    ModelGroup = "Media"
	ModelGroupOther    ModelGroup = "Other"
	ModelGroupGPT      ModelGroup = "GPT"
	ModelGroupClaude   ModelGroup = "Claude"
	ModelGroupGemini   ModelGroup = "Gemini"
	ModelGroupGemma    ModelGroup = "Gemma"
	ModelGroupGrok     ModelGroup = "Grok"
	ModelGroupCohere   ModelGroup = "Cohere"
	ModelGroupNova     ModelGroup = "Nova"
	ModelGroupQwen     ModelGroup = "Qwen"
	ModelGroupYi       ModelGroup = "Yi"
	ModelGroupDeepSeek ModelGroup = "DeepSeek"
	ModelGroupMistral  ModelGroup = "Mistral"
	ModelGroupLlama2   ModelGroup = "Llama2"
	ModelGroupLlama3   ModelGroup = "Llama3"
	ModelGroupLlama4   ModelGroup = "Llama4"
	ModelGroupPaLM     ModelGroup = "PaLM"
	ModelGroupRWKV     ModelGroup = "RWKV"
	ModelGroupQwen3    ModelGroup = "Qwen3"
)

// ReasoningEffort is a reasoning effort level, in descending effort order.
type ReasoningEffort string

const (
	ReasoningEffortMax     ReasoningEffort = "max"
	ReasoningEffortXHigh   ReasoningEffort = "xhigh"
	ReasoningEffortHigh    ReasoningEffort = "high"
	ReasoningEffortMedium  ReasoningEffort = "medium"
	ReasoningEffortLow     ReasoningEffort = "low"
	ReasoningEffortMinimal ReasoningEffort = "minimal"
	ReasoningEffortNone    ReasoningEffort = "none"
)

// ModelCategory is a use case category models can be filtered by.
type ModelCategory string

const (
	ModelCategoryProgramming  ModelCategory = "programming"
	ModelCategoryRoleplay     ModelCategory = "roleplay"
	ModelCategoryMarketing    ModelCategory = "marketing"
	ModelCategoryMarketingSEO ModelCategory = "marketing/seo"
	ModelCategoryTechnology   ModelCategory = "technology"
	ModelCategoryScience      ModelCategory = "science"
	ModelCategoryTranslation  ModelCategory = "translation"
	ModelCategoryLegal        ModelCategory = "legal"
	ModelCategoryFinance      ModelCategory = "finance"
	ModelCategoryHealth       ModelCategory = "health"
	ModelCategoryTrivia       ModelCategory = "trivia"
	ModelCategoryAcademia     ModelCategory = "academia"
)

// ModelSort is a server-side ordering of the models list. Models without a score
// for the chosen benchmark are placed last.
type ModelSort string

const (
	ModelSortMostPopular             ModelSort = "most-popular"
	ModelSortNewest                  ModelSort = "newest"
	ModelSortTopWeekly               ModelSort = "top-weekly"
	ModelSortPricingLowToHigh        ModelSort = "pricing-low-to-high"
	ModelSortPricingHighToLow        ModelSort = "pricing-high-to-low"
	ModelSortContextHighToLow        ModelSort = "context-high-to-low"
	ModelSortThroughputHighToLow     ModelSort = "throughput-high-to-low"
	ModelSortLatencyLowToHigh        ModelSort = "latency-low-to-high"
	ModelSortIntelligenceHighToLow   ModelSort = "intelligence-high-to-low"
	ModelSortCodingHighToLow         ModelSort = "coding-high-to-low"
	ModelSortAgenticHighToLow        ModelSort = "agentic-high-to-low"
	ModelSortDesignArenaELOHighToLow ModelSort = "design-arena-elo-high-to-low"
)

// ModelRegion is a data region models can be filtered by.
type ModelRegion string

const (
	ModelRegionEU ModelRegion = "eu"
)

// ListModelsOptions holds the optional query parameters of the models list
// endpoint. Zero values and nil pointers are omitted from the request.
type ListModelsOptions struct {
	Offset *int `url:"offset,omitempty"`
	Limit  *int `url:"limit,omitempty"`

	Category             ModelCategory    `url:"category,omitempty"`
	Sort                 ModelSort        `url:"sort,omitempty"`
	Search               string           `url:"search,omitempty"`
	Architecture         string           `url:"architecture,omitempty"`
	ModelAuthors         []string         `url:"model_authors,omitempty"`
	Providers            []string         `url:"providers,omitempty"`
	SupportedParameters  []Parameter      `url:"supported_parameters,omitempty"`
	InputModalities      []InputModality  `url:"input_modalities,omitempty"`
	OutputModalities     []OutputModality `url:"output_modalities,omitempty"`
	Context              *int             `url:"context,omitempty"`
	MinPrice             *float64         `url:"min_price,omitempty"`
	MaxPrice             *float64         `url:"max_price,omitempty"`
	MinOutputPrice       *float64         `url:"min_output_price,omitempty"`
	MaxOutputPrice       *float64         `url:"max_output_price,omitempty"`
	MinAgeDays           *int             `url:"min_age_days,omitempty"`
	MaxAgeDays           *int             `url:"max_age_days,omitempty"`
	MinIntelligenceIndex *float64         `url:"min_intelligence_index,omitempty"`
	MaxIntelligenceIndex *float64         `url:"max_intelligence_index,omitempty"`
	MinCodingIndex       *float64         `url:"min_coding_index,omitempty"`
	MaxCodingIndex       *float64         `url:"max_coding_index,omitempty"`
	MinAgenticIndex      *float64         `url:"min_agentic_index,omitempty"`
	MaxAgenticIndex      *float64         `url:"max_agentic_index,omitempty"`
	MinToolSuccessRate   *float64         `url:"min_tool_success_rate,omitempty"`
	MaxToolSuccessRate   *float64         `url:"max_tool_success_rate,omitempty"`
	Distillable          *bool            `url:"distillable,omitempty"`
	Region               ModelRegion      `url:"region,omitempty"`
	ZDR                  bool             `url:"zdr,omitempty"`
}

// ListUserModelsOptions holds the optional query parameters of the user models
// list endpoint.
type ListUserModelsOptions struct {
	Offset *int `url:"offset,omitempty"`
	Limit  *int `url:"limit,omitempty"`
}
