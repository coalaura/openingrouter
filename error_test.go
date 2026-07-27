package openingrouter

import (
	"testing"
)

func TestAsHttpError(t *testing.T) {
	client := tCreateClient(t)

	req, err := client.NewRequest("GET", "does/not/exist", nil)

	tAssertNil(t, err)

	_, err = client.Do(req)

	tAssertNotNil(t, err)
	tAssertEquals(t, err.Error(), "http: 404 Not Found")
}

func TestAsApiError(t *testing.T) {
	client := tCreateClient(t)

	req, err := client.NewRequest("GET", "models", nil)

	tAssertNil(t, err)

	req.URL.RawQuery = "context=0"

	_, err = client.Do(req)

	tAssertNotNil(t, err)
	tAssertEquals(t, err.Error(), "api error ZodError: Too small: expected number to be >=1")
}

func TestAsProviderError(t *testing.T) {
	client := tCreateClient(t)

	data := map[string]any{
		"model": "openai/gpt-5.6-luna",
		"messages": []any{
			map[string]any{
				"role": "user",
				"content": []any{
					map[string]any{
						"type": "text",
						"text": "hi",
					},
				},
			},
		},
		"response_format": map[string]any{
			"type": "json_object",
		},
	}

	req, err := client.NewRequest("POST", "chat/completions", data)

	tAssertNil(t, err)

	_, err = client.Do(req)

	tAssertNotNil(t, err)
	tAssertEquals(t, err.Error(), "provider error: Response input messages must contain the word 'json' in some form to use 'text.format' of type 'json_object'.")
}
