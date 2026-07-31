package openingrouter

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestAsHttpError(t *testing.T) {
	client := NewClient("none")

	req, err := client.NewRequest(context.Background(), "GET", "does/not/exist", nil)

	tAssertNil(t, err)

	_, err = client.Do(req)

	tAssertNotNil(t, err)
	tAssertEquals(t, err.Error(), "http: 404 Not Found")
}

func TestAsProviderError(t *testing.T) {
	err := AsOpenRouterError(tResponse(http.StatusBadRequest, `{
	"error": {
		"message": "Provider returned error",
		"code": 400,
		"metadata": {
			"raw": "{\n  \"error\": {\n    \"message\": \"Response input messages must contain the word 'json' in some form to use 'text.format' of type 'json_object'.\",\n    \"type\": \"invalid_request_error\",\n    \"param\": \"input\",\n    \"code\": null\n  }\n}",
			"provider_name": "Azure",
			"is_byok": false,
			"previous_errors": [
				{
					"code": 400,
					"message": "Provider returned error",
					"provider_name": "OpenAI",
					"raw": "{\n  \"error\": {\n    \"message\": \"Response input messages must contain the word 'json' in some form to use 'text.format' of type 'json_object'.\",\n    \"type\": \"invalid_request_error\",\n    \"param\": \"input\",\n    \"code\": null\n  }\n}"
				}
			]
		}
	},
	"user_id": "org_2xY1uKdHA6DpiQlHWHvh2jQIKf5"
}`, nil), nil)

	tAssertNotNil(t, err)

	provider, ok := errors.AsType[*ProviderError](err)

	tAssertEquals(t, ok, true)
	tAssertEquals(t, provider.Code, int64(400))
	tAssertEquals(t, provider.ProviderName, "Azure")
	tAssertEquals(t, provider.IsBYOK, false)

	tAssertEquals(t, provider.Type, ErrorType(""))
	tAssertEquals(t, provider.ProviderCode, "")
	tAssertEquals(t, provider.RetryAfter, time.Duration(0))
	tAssertEquals(t, provider.Retryable(), false)

	tAssertEquals(t, errors.Is(err, ErrInvalidRequest), true)
	tAssertEquals(t, errors.Is(err, ErrProviderUnavailable), false)

	tAssertEquals(t, err.Error(), "provider error: Response input messages must contain the word 'json' in some form to use 'text.format' of type 'json_object'.")

	tAssertEquals(t, len(provider.Previous), 1)
	tAssertEquals(t, provider.Previous[0].Code, int64(400))
	tAssertEquals(t, provider.Previous[0].ProviderName, "OpenAI")

	tAssertEquals(t, provider.Previous[0].Message, "Response input messages must contain the word 'json' in some form to use 'text.format' of type 'json_object'.")
}

func TestAsCreditsError(t *testing.T) {
	err := AsOpenRouterError(tResponse(http.StatusPaymentRequired, `{
	"error": {
		"message": "This request requires more credits, or fewer max_tokens. You requested up to 928596 tokens, but can only afford 377866. To increase, visit https://openrouter.ai/settings/credits and add more credits",
		"code": 402,
		"metadata": {
			"limit_source": "openrouter_credits",
			"remedy_hint": "Add credits at https://openrouter.ai/settings/credits, or lower max_tokens / prompt size to fit your remaining balance.",
			"provider_name": null,
			"previous_errors": [
				{
					"code": 402,
					"message": "This request requires more credits, or fewer max_tokens. You requested up to 928596 tokens, but can only afford 566799. To increase, visit https://openrouter.ai/settings/credits and add more credits"
				},
				{
					"code": 402,
					"message": "This request requires more credits, or fewer max_tokens. You requested up to 928596 tokens, but can only afford 566799. To increase, visit https://openrouter.ai/settings/credits and add more credits"
				}
			]
		}
	},
	"user_id": "org_2xY1uKdHA6DpiQlHWHvh2jQIKf5"
}`, nil), nil)

	tAssertNotNil(t, err)

	credits, ok := errors.AsType[*CreditsError](err)

	tAssertEquals(t, ok, true)
	tAssertEquals(t, credits.Code, int64(402))
	tAssertEquals(t, credits.LimitSource, "openrouter_credits")
	tAssertEquals(t, credits.RemedyHint, "Add credits at https://openrouter.ai/settings/credits, or lower max_tokens / prompt size to fit your remaining balance.")

	tAssertEquals(t, credits.Type, ErrorTypePaymentRequired)
	tAssertEquals(t, credits.ProviderCode, "")
	tAssertEquals(t, credits.RetryAfter, time.Duration(0))
	tAssertEquals(t, credits.Retryable(), false)

	tAssertEquals(t, errors.Is(err, ErrInsufficientCredits), true)
	tAssertEquals(t, errors.Is(err, ErrInvalidRequest), false)
	tAssertEquals(t, errors.Is(err, ErrRateLimited), false)

	tAssertEquals(t, err.Error(), "insufficient credits: This request requires more credits, or fewer max_tokens. You requested up to 928596 tokens, but can only afford 377866. To increase, visit https://openrouter.ai/settings/credits and add more credits")

	tAssertEquals(t, len(credits.Previous), 2)
	tAssertEquals(t, credits.Previous[0].Code, int64(402))
	tAssertEquals(t, credits.Previous[0].Message, "This request requires more credits, or fewer max_tokens. You requested up to 928596 tokens, but can only afford 566799. To increase, visit https://openrouter.ai/settings/credits and add more credits")
	tAssertEquals(t, credits.Previous[1].Message, credits.Previous[0].Message)
}

func TestAsOpenRouterErrorShapes(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		header  http.Header
		body    string
		want    string
		wantIs  error
		inspect func(t *testing.T, err error)
	}{
		{
			name:   "typed provider error with retry-after",
			status: http.StatusTooManyRequests,
			header: http.Header{"Retry-After": []string{"60"}},
			body:   `{"error":{"code":429,"message":"Provider returned error","metadata":{"raw":"{\"error\":{\"message\":\"Rate limit exceeded\"}}","provider_name":"OpenAI","error_type":"rate_limit_exceeded","provider_code":"rate_limited"}}}`,
			want:   "provider error: Rate limit exceeded",
			wantIs: ErrRateLimited,
			inspect: func(t *testing.T, err error) {
				provider, ok := errors.AsType[*ProviderError](err)

				tAssertEquals(t, ok, true)
				tAssertEquals(t, provider.ProviderCode, "rate_limited")
				tAssertEquals(t, provider.RetryAfter, 60*time.Second)
				tAssertEquals(t, provider.Retryable(), true)
			},
		},
		{
			name:   "moderation",
			status: http.StatusForbidden,
			body:   `{"error":{"code":403,"message":"Input flagged","metadata":{"reasons":["violence","hate"],"flagged_input":"...","provider_name":"Anthropic","model_slug":"anthropic/claude-haiku-4.5"}}}`,
			want:   "moderation error: Input flagged (violence, hate)",
			wantIs: ErrModerated,
			inspect: func(t *testing.T, err error) {
				moderation, ok := errors.AsType[*ModerationError](err)

				tAssertEquals(t, ok, true)
				tAssertEquals(t, moderation.Type, ErrorTypeContentPolicy)
				tAssertEquals(t, moderation.ModelSlug, "anthropic/claude-haiku-4.5")
			},
		},
		{
			name:   "guardrail falls back to unmodeled metadata",
			status: http.StatusForbidden,
			body:   `{"error":{"code":403,"message":"Request blocked: prompt injection patterns detected","metadata":{"patterns":["ignore all previous instructions"]}}}`,
			want:   "openrouter code 403: Request blocked: prompt injection patterns detected",
			wantIs: ErrForbidden,
			inspect: func(t *testing.T, err error) {
				generic, ok := errors.AsType[*OpenRouterError](err)

				tAssertEquals(t, ok, true)
				tAssertEquals(t, string(generic.Metadata), `{"patterns":["ignore all previous instructions"]}`)
			},
		},
		{
			name:   "named api error",
			status: http.StatusBadRequest,
			body:   `{"success":false,"error":{"name":"ZodError","message": "[\n  {\n    \"code\": \"invalid_value\",\n    \"values\": [\n      \"programming\",\n      \"roleplay\",\n      \"marketing\",\n      \"marketing/seo\",\n      \"technology\",\n      \"science\",\n      \"translation\",\n      \"legal\",\n      \"finance\",\n      \"health\",\n      \"trivia\",\n      \"academia\"\n    ],\n    \"path\": [\n      \"category\"\n    ],\n    \"message\": \"Invalid option: expected one of \\\"programming\\\"|\\\"roleplay\\\"|\\\"marketing\\\"|\\\"marketing/seo\\\"|\\\"technology\\\"|\\\"science\\\"|\\\"translation\\\"|\\\"legal\\\"|\\\"finance\\\"|\\\"health\\\"|\\\"trivia\\\"|\\\"academia\\\"\"\n  }\n]"}}`,
			want:   `api error ZodError: Invalid option: expected one of "programming"|"roleplay"|"marketing"|"marketing/seo"|"technology"|"science"|"translation"|"legal"|"finance"|"health"|"trivia"|"academia"`,
			wantIs: ErrInvalidRequest,
			inspect: func(t *testing.T, err error) {
				apiErr, ok := errors.AsType[*ApiError](err)

				tAssertEquals(t, ok, true)
				tAssertEquals(t, apiErr.Name, "ZodError")

				tAssertEquals(t, apiErr.Code, int64(400))
			},
		},
		{
			name:   "context length via error_type on a masked status",
			status: http.StatusOK,
			body:   `{"error":{"code":400,"message":"too long","metadata":{"error_type":"context_length_exceeded"}}}`,
			want:   "openrouter code 400: too long",
			wantIs: ErrContextLength,
		},
		{
			name:   "undecodable body",
			status: http.StatusNotFound,
			body:   "<html>nope</html>",
			want:   "http: 404 Not Found",
			wantIs: ErrNotFound,
		},
		{
			name:   "empty error object",
			status: http.StatusBadGateway,
			body:   `{"error":{}}`,
			want:   "http: 502 Bad Gateway",
			wantIs: ErrProviderUnavailable,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := AsOpenRouterError(tResponse(test.status, test.body, test.header), nil)

			tAssertNotNil(t, err)
			tAssertEquals(t, err.Error(), test.want)
			tAssertEquals(t, errors.Is(err, test.wantIs), true)

			if test.inspect != nil {
				test.inspect(t, err)
			}
		})
	}
}

func TestAsOpenRouterErrorPassesThroughTransportError(t *testing.T) {
	want := errors.New("dial tcp: connection refused")

	tAssertEquals(t, AsOpenRouterError(nil, want), want)
}
