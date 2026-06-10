package openingrouter

import (
	"context"
	"encoding/json"
	"net/http"
)

// Root response type
type FrontendModelsResponse struct {
	Data []FrontendModel `json:"data"`
}

// FrontendModel represents a model in the OpenRouter frontend API
type FrontendModel struct {
	Slug                  string                   `json:"slug"`
	HFSlug                string                   `json:"hf_slug"`
	UpdatedAt             FlexibleTime             `json:"updated_at"`
	CreatedAt             FlexibleTime             `json:"created_at"`
	HFUpdatedAt           *FlexibleTime            `json:"hf_updated_at"`
	Name                  string                   `json:"name"`
	ShortName             string                   `json:"short_name"`
	Author                string                   `json:"author"`
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
	WarningMessage        string                   `json:"warning_message"`
	PromotionMessage      string                   `json:"promotion_message"`
	RoutingErrorMessage   string                   `json:"routing_error_message"`
	Permaslug             string                   `json:"permaslug"`
	SupportsReasoning     bool                     `json:"supports_reasoning"`
	ReasoningConfig       *FrontendReasoningConfig `json:"reasoning_config"`
	Features              *FrontendFeatures        `json:"features"`
	DefaultParameters     *FrontendDefaultParams   `json:"default_parameters"`
	DefaultOrder          []string                 `json:"default_order"`
	QuickStartExampleType string                   `json:"quick_start_example_type"`
	IsTrainableText       *bool                    `json:"is_trainable_text"`
	IsTrainableImage      *bool                    `json:"is_trainable_image"`
	Endpoint              *FrontendEndpoint        `json:"endpoint"`
}

// FrontendReasoningConfig represents reasoning configuration
type FrontendReasoningConfig struct {
	StartToken                *string  `json:"start_token"`
	EndToken                  *string  `json:"end_token"`
	IsMandatoryReasoning      *bool    `json:"is_mandatory_reasoning"`
	SupportsReasoningEffort   *bool    `json:"supports_reasoning_effort"`
	SupportedReasoningEfforts []string `json:"supported_reasoning_efforts"`
}

// FrontendFeatures represents model features
type FrontendFeatures struct {
	ReasoningConfig     *FrontendReasoningConfig `json:"reasoning_config"`
	ChatTemplateConfig  map[string]any           `json:"chat_template_config"`
	SupportsFileURLs    *bool                    `json:"supports_file_urls,omitempty"`
	SupportsBase64Video *bool                    `json:"supports_base64_video_input,omitempty"`
	SupportsVideoURLs   *bool                    `json:"supports_video_urls,omitempty"`
	SupportsToolChoice  *FrontendToolChoice      `json:"supports_tool_choice,omitempty"`
	SupportsInputAudio  *bool                    `json:"supports_input_audio,omitempty"`
	SupportsNativeWeb   *bool                    `json:"supports_native_web_search,omitempty"`
	SupportsMultipart   *bool                    `json:"supports_multipart,omitempty"`
}

// FrontendToolChoice represents tool choice support
type FrontendToolChoice struct {
	LiteralNone     bool `json:"literal_none"`
	LiteralAuto     bool `json:"literal_auto"`
	LiteralRequired bool `json:"literal_required"`
	TypeFunction    bool `json:"type_function"`
}

// FrontendDefaultParams represents default model parameters
type FrontendDefaultParams struct {
	Temperature      *float64 `json:"temperature"`
	TopP             *float64 `json:"top_p"`
	FrequencyPenalty *float64 `json:"frequency_penalty"`
}

// FrontendEndpoint represents an endpoint configuration
type FrontendEndpoint struct {
	ID                    string                    `json:"id"`
	Name                  string                    `json:"name"`
	ContextLength         int                       `json:"context_length"`
	Model                 *FrontendModel            `json:"model"`
	ModelVariantSlug      string                    `json:"model_variant_slug"`
	ModelVariantPermaslug string                    `json:"model_variant_permaslug"`
	AdapterName           string                    `json:"adapter_name"`
	ProviderName          string                    `json:"provider_name"`
	ProviderInfo          *FrontendProviderInfo     `json:"provider_info"`
	ProviderDisplayName   string                    `json:"provider_display_name"`
	ProviderSlug          string                    `json:"provider_slug"`
	ProviderModelID       string                    `json:"provider_model_id"`
	Quantization          string                    `json:"quantization"`
	Variant               string                    `json:"variant"`
	IsFree                bool                      `json:"is_free"`
	CanAbort              bool                      `json:"can_abort"`
	MaxPromptTokens       *int                      `json:"max_prompt_tokens"`
	MaxCompletionTokens   int                       `json:"max_completion_tokens"`
	MaxTokensPerImage     *int                      `json:"max_tokens_per_image"`
	SupportedParameters   []string                  `json:"supported_parameters"`
	IsBYOK                bool                      `json:"is_byok"`
	ModerationRequired    bool                      `json:"moderation_required"`
	DataPolicy            *FrontendDataPolicy       `json:"data_policy"`
	Pricing               *FrontendPricing          `json:"pricing"`
	VariablePricings      []FrontendVariablePricing `json:"variable_pricings"`
	LineItems             []FrontendLineItem        `json:"line_items"`
	PricingJSON           map[string]any            `json:"pricing_json"`
	PricingVersionID      string                    `json:"pricing_version_id"`
	IsHidden              bool                      `json:"is_hidden"`
	IsDeranked            bool                      `json:"is_deranked"`
	IsDisabled            bool                      `json:"is_disabled"`
	SupportsToolParams    bool                      `json:"supports_tool_parameters"`
	SupportsReasoning     bool                      `json:"supports_reasoning"`
	SupportsMultipart     bool                      `json:"supports_multipart"`
	LimitRPM              *int                      `json:"limit_rpm"`
	LimitRPD              *int                      `json:"limit_rpd"`
	LimitRPMCF            *int                      `json:"limit_rpm_cf"`
	HasCompletions        bool                      `json:"has_completions"`
	HasChatCompletions    bool                      `json:"has_chat_completions"`
	Features              *FrontendFeatures         `json:"features"`
	ProviderRegion        *string                   `json:"provider_region"`
	DeprecationDate       *FlexibleTime             `json:"deprecation_date"`
}

// FrontendProviderInfo represents provider information
type FrontendProviderInfo struct {
	Name                  string                            `json:"name"`
	DisplayName           string                            `json:"displayName"`
	Slug                  string                            `json:"slug"`
	BaseURL               string                            `json:"baseUrl"`
	DataPolicy            *FrontendDataPolicy               `json:"dataPolicy"`
	Headquarters          string                            `json:"headquarters"`
	Datacenters           []string                          `json:"datacenters,omitempty"`
	RegionOverrides       map[string]FrontendRegionOverride `json:"regionOverrides"`
	HasChatCompletions    bool                              `json:"hasChatCompletions"`
	HasCompletions        bool                              `json:"hasCompletions"`
	IsAbortable           bool                              `json:"isAbortable"`
	ModerationRequired    bool                              `json:"moderationRequired"`
	Editors               []string                          `json:"editors"`
	Owners                []string                          `json:"owners"`
	AdapterName           string                            `json:"adapterName"`
	IsMultipartSupported  bool                              `json:"isMultipartSupported"`
	StatusPageURL         string                            `json:"statusPageUrl"`
	BYOKEnabled           bool                              `json:"byokEnabled"`
	Icon                  *FrontendIcon                     `json:"icon"`
	IgnoredProviderModels []string                          `json:"ignoredProviderModels"`
	SendClientIP          bool                              `json:"sendClientIp"`
	PricingStrategy       string                            `json:"pricingStrategy"`
}

// FrontendRegionOverride represents region-specific overrides
type FrontendRegionOverride struct {
	BaseURL string `json:"baseUrl"`
}

// FrontendIcon represents a provider icon
type FrontendIcon struct {
	URL string `json:"url"`
}

// FrontendDataPolicy represents data policy configuration
type FrontendDataPolicy struct {
	Training           bool   `json:"training"`
	TrainingOpenRouter bool   `json:"trainingOpenRouter"`
	RetainsPrompts     bool   `json:"retainsPrompts"`
	RetentionDays      *int   `json:"retentionDays,omitempty"`
	CanPublish         bool   `json:"canPublish"`
	TermsOfServiceURL  string `json:"termsOfServiceURL"`
	PrivacyPolicyURL   string `json:"privacyPolicyURL"`
	RequiresUserIDs    bool   `json:"requiresUserIDs"`
}

// FrontendPricing represents pricing information
type FrontendPricing struct {
	Prompt            StringifiedNumber  `json:"prompt"`
	Completion        StringifiedNumber  `json:"completion"`
	Image             StringifiedNumber  `json:"image,omitempty"`
	Audio             StringifiedNumber  `json:"audio,omitempty"`
	InputAudioCache   StringifiedNumber  `json:"input_audio_cache,omitempty"`
	InputCacheRead    StringifiedNumber  `json:"input_cache_read,omitempty"`
	InputCacheWrite   StringifiedNumber  `json:"input_cache_write,omitempty"`
	InternalReasoning StringifiedNumber  `json:"internal_reasoning,omitempty"`
	WebSearch         StringifiedNumber  `json:"web_search,omitempty"`
	Discount          float64            `json:"discount"`
	LineItems         []FrontendLineItem `json:"line_items"`
}

// FrontendVariablePricing represents variable pricing tiers
type FrontendVariablePricing struct {
	Type            string            `json:"type"`
	Threshold       any               `json:"threshold"`
	Prompt          StringifiedNumber `json:"prompt"`
	Completions     StringifiedNumber `json:"completions"`
	InputCacheRead  StringifiedNumber `json:"input_cache_read,omitempty"`
	InputCacheWrite StringifiedNumber `json:"input_cache_write,omitempty"`
}

// FrontendLineItem represents a pricing line item
type FrontendLineItem struct {
	Type  string `json:"type"`
	Value string `json:"value"`
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
