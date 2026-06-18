package openingrouter

import (
	"context"
	"encoding/json"
	"net/http"
)

// FrontendModelsResponse is the root response for the catalog models endpoint.
type FrontendModelsResponse struct {
	Data []FrontendModel `json:"data"`
}

// FrontendModel represents a model in the OpenRouter frontend API.
type FrontendModel struct {
	Slug                  string                   `json:"slug"`
	HFSlug                *string                  `json:"hf_slug"`
	UpdatedAt             FlexibleTime             `json:"updated_at"`
	CreatedAt             FlexibleTime             `json:"created_at"`
	HFUpdatedAt           *FlexibleTime            `json:"hf_updated_at"`
	Name                  string                   `json:"name"`
	ShortName             string                   `json:"short_name"`
	Author                string                   `json:"author"`
	AuthorDisplayName     string                   `json:"author_display_name"`
	Description           string                   `json:"description"`
	ModelVersionGroupID   *string                  `json:"model_version_group_id"`
	ContextLength         int                      `json:"context_length"`
	InputModalities       []string                 `json:"input_modalities"`
	OutputModalities      []string                 `json:"output_modalities"`
	HasTextOutput         bool                     `json:"has_text_output"`
	Group                 string                   `json:"group"`
	InstructType          *string                  `json:"instruct_type"`
	DefaultSystem         *string                  `json:"default_system"`
	DefaultStops          []string                 `json:"default_stops"`
	Hidden                bool                     `json:"hidden"`
	Router                *string                  `json:"router"`
	WarningMessage        *string                  `json:"warning_message"`
	PromotionMessage      *string                  `json:"promotion_message"`
	RoutingErrorMessage   *string                  `json:"routing_error_message"`
	IsPrivate             bool                     `json:"is_private"`
	Permaslug             string                   `json:"permaslug"`
	SupportsReasoning     bool                     `json:"supports_reasoning"`
	ReasoningConfig       *FrontendReasoningConfig `json:"reasoning_config"`
	Features              *FrontendFeatures        `json:"features"`
	DefaultParameters     *FrontendDefaultParams   `json:"default_parameters"`
	DefaultOrder          []string                 `json:"default_order"`
	QuickStartExampleType string                   `json:"quick_start_example_type"`
	IsTrainableText       *bool                    `json:"is_trainable_text"`
	IsTrainableImage      *bool                    `json:"is_trainable_image"`
	KnowledgeCutoff       *string                  `json:"knowledge_cutoff"`
	LimitRPM              *int                     `json:"limit_rpm"`
	LimitRPD              *int                     `json:"limit_rpd"`
	SupportedTTSVoices    []string                 `json:"supported_tts_voices"`
	Endpoint              *FrontendEndpoint        `json:"endpoint,omitempty"`
}

// FrontendReasoningConfig represents reasoning configuration for a model or endpoint.
type FrontendReasoningConfig struct {
	StartToken                *string  `json:"start_token"`
	EndToken                  *string  `json:"end_token"`
	IsMandatoryReasoning      *bool    `json:"is_mandatory_reasoning,omitempty"`
	SupportsReasoningEffort   *bool    `json:"supports_reasoning_effort,omitempty"`
	SupportedReasoningEfforts []string `json:"supported_reasoning_efforts"`
	DefaultReasoningEffort    *string  `json:"default_reasoning_effort,omitempty"`
	DefaultReasoningEnabled   *bool    `json:"default_reasoning_enabled,omitempty"`
	ReasoningReturnMechanism  *string  `json:"reasoning_return_mechanism,omitempty"`
}

// FrontendFeatures represents model and endpoint feature support.
type FrontendFeatures struct {
	ReasoningConfig          *FrontendReasoningConfig `json:"reasoning_config,omitempty"`
	ChatTemplateConfig       map[string]any           `json:"chat_template_config,omitempty"`
	ReasoningReturnMechanism *string                  `json:"reasoning_return_mechanism,omitempty"`
	SupportsFileURLs         *bool                    `json:"supports_file_urls,omitempty"`
	SupportsBase64Video      *bool                    `json:"supports_base64_video_input,omitempty"`
	SupportsVideoURLs        *bool                    `json:"supports_video_urls,omitempty"`
	SupportsToolChoice       *FrontendToolChoice      `json:"supports_tool_choice,omitempty"`
	SupportsInputAudio       *bool                    `json:"supports_input_audio,omitempty"`
	SupportsNativeWeb        *bool                    `json:"supports_native_web_search,omitempty"`
	SupportsMultipart        *bool                    `json:"supports_multipart,omitempty"`
}

// FrontendSupportedVideoParameters represents video generation parameter support for an endpoint.
type FrontendSupportedVideoParameters struct {
	SupportedResolutions  []string `json:"supported_resolutions"`
	SupportedAspectRatios []string `json:"supported_aspect_ratios"`
	SupportedSizes        []string `json:"supported_sizes"`
	SupportedDurations    []int    `json:"supported_durations"`
	SupportedFrameImages  []string `json:"supported_frame_images"`
	GenerateAudio         *bool    `json:"generate_audio"`
	Seed                  *bool    `json:"seed"`
}

// FrontendToolChoice represents which tool_choice values an endpoint supports.
type FrontendToolChoice struct {
	LiteralNone     bool `json:"literal_none"`
	LiteralAuto     bool `json:"literal_auto"`
	LiteralRequired bool `json:"literal_required"`
	TypeFunction    bool `json:"type_function"`
}

// FrontendDefaultParams represents default sampling parameters for a model.
type FrontendDefaultParams struct {
	Temperature       *float64 `json:"temperature"`
	TopP              *float64 `json:"top_p"`
	TopK              *float64 `json:"top_k"`
	FrequencyPenalty  *float64 `json:"frequency_penalty"`
	PresencePenalty   *float64 `json:"presence_penalty"`
	RepetitionPenalty *float64 `json:"repetition_penalty"`
}

// FrontendEndpoint represents a provider endpoint serving a model.
type FrontendEndpoint struct {
	ID                           string                            `json:"id"`
	Name                         string                            `json:"name"`
	ContextLength                int                               `json:"context_length"`
	Model                        *FrontendModel                    `json:"model"`
	ModelVariantSlug             string                            `json:"model_variant_slug"`
	ModelVariantPermaslug        string                            `json:"model_variant_permaslug"`
	AdapterName                  string                            `json:"adapter_name"`
	ProviderName                 string                            `json:"provider_name"`
	ProviderInfo                 *FrontendProviderInfo             `json:"provider_info"`
	ProviderDisplayName          string                            `json:"provider_display_name"`
	ProviderSlug                 string                            `json:"provider_slug"`
	ProviderModelID              string                            `json:"provider_model_id"`
	Quantization                 string                            `json:"quantization"`
	Variant                      string                            `json:"variant"`
	IsFree                       bool                              `json:"is_free"`
	CanAbort                     bool                              `json:"can_abort"`
	MaxPromptTokens              *int                              `json:"max_prompt_tokens"`
	MaxCompletionTokens          int                               `json:"max_completion_tokens"`
	MaxTokensPerImage            *int                              `json:"max_tokens_per_image"`
	SupportedParameters          []string                          `json:"supported_parameters"`
	ExcludedParameters           []string                          `json:"excluded_parameters"`
	IsBYOK                       bool                              `json:"is_byok"`
	ModerationRequired           bool                              `json:"moderation_required"`
	DataPolicy                   *FrontendDataPolicy               `json:"data_policy"`
	Pricing                      *FrontendPricing                  `json:"pricing"`
	DisplayPricing               []FrontendDisplayPricing          `json:"display_pricing"`
	VariablePricings             []FrontendVariablePricing         `json:"variable_pricings,omitempty"`
	LineItems                    []FrontendLineItem                `json:"line_items,omitempty"`
	PricingJSON                  map[string]any                    `json:"pricing_json"`
	PricingVersionID             string                            `json:"pricing_version_id"`
	IsHidden                     bool                              `json:"is_hidden"`
	IsPrivate                    bool                              `json:"is_private"`
	IsDeranked                   bool                              `json:"is_deranked"`
	IsDisabled                   bool                              `json:"is_disabled"`
	SupportsToolParams           bool                              `json:"supports_tool_parameters"`
	SupportsReasoning            bool                              `json:"supports_reasoning"`
	SupportsMultipart            bool                              `json:"supports_multipart"`
	LimitRPM                     *int                              `json:"limit_rpm"`
	LimitRPD                     *int                              `json:"limit_rpd"`
	LimitRPMCF                   *int                              `json:"limit_rpm_cf"`
	HasCompletions               bool                              `json:"has_completions"`
	HasChatCompletions           bool                              `json:"has_chat_completions"`
	Features                     *FrontendFeatures                 `json:"features"`
	SupportedVideoParameters     *FrontendSupportedVideoParameters `json:"supported_video_parameters"`
	ProviderRegion               *string                           `json:"provider_region"`
	DeprecationDate              *FlexibleTime                     `json:"deprecation_date"`
	AllowedPassthroughParameters []string                          `json:"allowed_passthrough_parameters"`
	CapacityTPM                  *float64                          `json:"capacity_tpm"`
	CreatedAt                    FlexibleTime                      `json:"created_at"`
	RoutingHeuristics            *FrontendRoutingHeuristics        `json:"routing_heuristics"`
	Status                       int                               `json:"status"`
	StatusHeuristics             *FrontendStatusHeuristics         `json:"status_heuristics"`
	StatusHeuristics5m           *FrontendStatusHeuristics         `json:"status_heuristics_5m"`
	StatusHeuristics1d           *FrontendStatusHeuristics         `json:"status_heuristics_1d"`
	Fortuna                      *FrontendFortuna                  `json:"fortuna"`
}

// FrontendProviderInfo represents information about an inference provider.
type FrontendProviderInfo struct {
	Name                  string                            `json:"name"`
	DisplayName           string                            `json:"displayName"`
	Slug                  string                            `json:"slug"`
	BaseURL               string                            `json:"baseUrl"`
	DataPolicy            *FrontendDataPolicy               `json:"dataPolicy"`
	Headquarters          string                            `json:"headquarters"`
	Datacenters           []string                          `json:"datacenters,omitempty"`
	RegionOverrides       map[string]FrontendRegionOverride `json:"regionOverrides,omitempty"`
	HasChatCompletions    bool                              `json:"hasChatCompletions"`
	HasCompletions        bool                              `json:"hasCompletions"`
	IsAbortable           bool                              `json:"isAbortable"`
	ModerationRequired    bool                              `json:"moderationRequired"`
	Editors               []string                          `json:"editors"`
	Owners                []string                          `json:"owners"`
	AdapterName           string                            `json:"adapterName"`
	IsMultipartSupported  bool                              `json:"isMultipartSupported,omitempty"`
	StatusPageURL         *string                           `json:"statusPageUrl"`
	BYOKEnabled           bool                              `json:"byokEnabled"`
	Icon                  *FrontendIcon                     `json:"icon"`
	IgnoredProviderModels []string                          `json:"ignoredProviderModels,omitempty"`
	SendClientIP          bool                              `json:"sendClientIp"`
	PricingStrategy       string                            `json:"pricingStrategy"`
}

// FrontendRegionOverride represents region-specific provider overrides.
type FrontendRegionOverride struct {
	BaseURL string `json:"baseUrl"`
}

// FrontendIcon represents a provider icon.
type FrontendIcon struct {
	URL string `json:"url"`
}

// FrontendDataPolicy represents a provider's data handling policy.
type FrontendDataPolicy struct {
	Training           bool   `json:"training"`
	TrainingOpenRouter bool   `json:"trainingOpenRouter"`
	RetainsPrompts     bool   `json:"retainsPrompts"`
	RetentionDays      *int   `json:"retentionDays,omitempty"`
	CanPublish         bool   `json:"canPublish"`
	TermsOfServiceURL  string `json:"termsOfServiceURL"`
	PrivacyPolicyURL   string `json:"privacyPolicyURL"`
	RequiresUserIDs    bool   `json:"requiresUserIDs,omitempty"`
}

// FrontendPricing represents endpoint pricing information.
type FrontendPricing struct {
	Prompt            StringifiedNumber        `json:"prompt"`
	Completion        StringifiedNumber        `json:"completion"`
	Image             StringifiedNumber        `json:"image,omitempty"`
	ImageOutput       StringifiedNumber        `json:"image_output,omitempty"`
	Audio             StringifiedNumber        `json:"audio,omitempty"`
	InputAudioCache   StringifiedNumber        `json:"input_audio_cache,omitempty"`
	InputCacheRead    StringifiedNumber        `json:"input_cache_read,omitempty"`
	InputCacheWrite   StringifiedNumber        `json:"input_cache_write,omitempty"`
	InternalReasoning StringifiedNumber        `json:"internal_reasoning,omitempty"`
	WebSearch         StringifiedNumber        `json:"web_search,omitempty"`
	Discount          float64                  `json:"discount"`
	DisplayPricing    []FrontendDisplayPricing `json:"display_pricing,omitempty"`
	LineItems         []FrontendLineItem       `json:"line_items,omitempty"`
}

// FrontendDisplayPricing represents a human-readable pricing line item.
type FrontendDisplayPricing struct {
	Kind              string `json:"kind"`
	SKULabel          string `json:"sku_label"`
	Price             string `json:"price"`
	DisplayMultiplier int    `json:"displayMultiplier"`
	UnitLabel         string `json:"unitLabel"`
}

// FrontendVariablePricing represents variable pricing tiers.
type FrontendVariablePricing struct {
	Type            string            `json:"type"`
	Threshold       any               `json:"threshold"`
	Prompt          StringifiedNumber `json:"prompt"`
	Completions     StringifiedNumber `json:"completions"`
	InputCacheRead  StringifiedNumber `json:"input_cache_read,omitempty"`
	InputCacheWrite StringifiedNumber `json:"input_cache_write,omitempty"`
}

// FrontendLineItem represents a pricing line item.
type FrontendLineItem struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// FrontendRoutingHeuristics represents throughput, latency and routing statistics for an endpoint.
type FrontendRoutingHeuristics struct {
	RequestCount                       int                   `json:"request_count"`
	P50Throughput                      float64               `json:"p50_throughput"`
	P50Latency                         float64               `json:"p50_latency"`
	RequestCount30Minutes              int                   `json:"request_count_30_minutes"`
	P50Throughput30Minutes             float64               `json:"p50_throughput_30_minutes"`
	P75Throughput30Minutes             float64               `json:"p75_throughput_30_minutes"`
	P90Throughput30Minutes             float64               `json:"p90_throughput_30_minutes"`
	P99Throughput30Minutes             float64               `json:"p99_throughput_30_minutes"`
	P50Latency30Minutes                float64               `json:"p50_latency_30_minutes"`
	P75Latency30Minutes                float64               `json:"p75_latency_30_minutes"`
	P90Latency30Minutes                float64               `json:"p90_latency_30_minutes"`
	P99Latency30Minutes                float64               `json:"p99_latency_30_minutes"`
	EffectivePromptPrice               float64               `json:"effective_prompt_price"`
	EffectiveCompletionPrice           float64               `json:"effective_completion_price"`
	RequestCount5Minutes               int                   `json:"request_count_5_minutes"`
	P50Throughput5Minutes              float64               `json:"p50_throughput_5_minutes"`
	P75Throughput5Minutes              float64               `json:"p75_throughput_5_minutes"`
	P90Throughput5Minutes              float64               `json:"p90_throughput_5_minutes"`
	P99Throughput5Minutes              float64               `json:"p99_throughput_5_minutes"`
	P50Latency5Minutes                 float64               `json:"p50_latency_5_minutes"`
	P75Latency5Minutes                 float64               `json:"p75_latency_5_minutes"`
	P90Latency5Minutes                 float64               `json:"p90_latency_5_minutes"`
	P99Latency5Minutes                 float64               `json:"p99_latency_5_minutes"`
	P50Throughput2Hours                float64               `json:"p50_throughput_2_hours"`
	ToolFinishReasonRequestSuccessRate *float64              `json:"tool_finish_reason_request_success_rate,omitempty"`
	ToolCallsFinishReasonRequestCount  *int                  `json:"tool_calls_finish_reason_request_count,omitempty"`
	DerankResult                       *FrontendDerankResult `json:"derank_result,omitempty"`
}

// FrontendDerankResult represents whether and why an endpoint is deranked.
type FrontendDerankResult struct {
	IsDeranked    bool                   `json:"is_deranked"`
	DerankReasons []FrontendDerankReason `json:"derank_reasons"`
}

// FrontendDerankReason represents a single reason an endpoint was deranked.
type FrontendDerankReason struct {
	Signal             string  `json:"signal"`
	StddevsBelowMedian float64 `json:"stddevs_below_median"`
	Threshold          float64 `json:"threshold"`
	Value              float64 `json:"value"`
}

// FrontendStatusHeuristics represents request outcome counts over a time window.
type FrontendStatusHeuristics struct {
	Success         int `json:"success"`
	DerankableError int `json:"derankableError"`
	RateLimited     int `json:"rateLimited"`
}

// FrontendFortuna represents capacity and load-balancing scoring for an endpoint.
type FrontendFortuna struct {
	BetaAlpha          float64 `json:"beta_alpha"`
	BetaBeta           float64 `json:"beta_beta"`
	CapacityScore      float64 `json:"capacity_score"`
	CapacityCeilingRPM float64 `json:"capacity_ceiling_rpm"`
	RecentPeakRPM      float64 `json:"recent_peak_rpm"`
}

func ListFrontendModels(ctx context.Context) ([]FrontendModel, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://openrouter.ai/api/frontend/v1/catalog/models", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result FrontendModelsResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}
