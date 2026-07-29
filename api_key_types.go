package openingrouter

// ApiKeyInfo represents information on the API key associated with the current
// authentication session.
type ApiKeyInfo struct {
	Label              string            `json:"label"`
	Usage              float64           `json:"usage"`
	UsageDaily         float64           `json:"usage_daily"`
	UsageWeekly        float64           `json:"usage_weekly"`
	UsageMonthly       float64           `json:"usage_monthly"`
	BYOKUsage          float64           `json:"byok_usage"`
	BYOKUsageDaily     float64           `json:"byok_usage_daily"`
	BYOKUsageWeekly    float64           `json:"byok_usage_weekly"`
	BYOKUsageMonthly   float64           `json:"byok_usage_monthly"`
	Limit              *float64          `json:"limit"`
	LimitRemaining     *float64          `json:"limit_remaining"`
	LimitReset         *ApiKeyLimitReset `json:"limit_reset"`
	IncludeBYOKInLimit bool              `json:"include_byok_in_limit"`
	IsFreeTier         bool              `json:"is_free_tier"`
	IsManagementKey    bool              `json:"is_management_key"`
	IsProvisioningKey  bool              `json:"is_provisioning_key"`
	CreatorUserID      *string           `json:"creator_user_id"`
	ExpiresAt          *FlexibleTime     `json:"expires_at"`
	RateLimit          ApiKeyRateLimit   `json:"rate_limit"`
}

// ApiKeyRateLimit represents legacy rate limit information about an API key.
// Deprecated: this field is safe to ignore.
type ApiKeyRateLimit struct {
	Requests int    `json:"requests"`
	Interval string `json:"interval"`
	Note     string `json:"note"`
}

// ApiKeyLimitReset is the type of limit reset schedule for an API key.
type ApiKeyLimitReset string

const (
	ApiKeyLimitResetDaily   ApiKeyLimitReset = "daily"
	ApiKeyLimitResetWeekly  ApiKeyLimitReset = "weekly"
	ApiKeyLimitResetMonthly ApiKeyLimitReset = "monthly"
)
