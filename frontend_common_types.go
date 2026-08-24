package openingrouter

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
