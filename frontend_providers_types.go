package openingrouter

// FrontendProvider represents an inference provider in the frontend API.
type FrontendProvider struct {
	Name               string              `json:"name"`
	DisplayName        string              `json:"displayName"`
	Slug               string              `json:"slug"`
	AdapterName        string              `json:"adapterName"`
	BaseURL            string              `json:"baseUrl"`
	DataPolicy         *FrontendDataPolicy `json:"dataPolicy"`
	Headquarters       string              `json:"headquarters,omitempty"`
	Datacenters        []string            `json:"datacenters,omitempty"`
	HasChatCompletions bool                `json:"hasChatCompletions"`
	HasCompletions     bool                `json:"hasCompletions"`
	IsAbortable        bool                `json:"isAbortable"`
	ModerationRequired bool                `json:"moderationRequired"`
	StatusPageURL      *string             `json:"statusPageUrl"`
	BYOKEnabled        bool                `json:"byokEnabled"`
	Icon               *FrontendIcon       `json:"icon"`
	SendClientIP       bool                `json:"sendClientIp"`
	PricingStrategy    string              `json:"pricingStrategy"`
}
