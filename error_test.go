package openingrouter

import (
	"context"
	"testing"
)

func TestAsHttpError(t *testing.T) {
	client := tCreateClient(t)

	req, err := client.NewRequest(context.Background(), "GET", "does/not/exist", nil)

	tAssertNil(t, err)

	_, err = client.Do(req)

	tAssertNotNil(t, err)
	tAssertEquals(t, err.Error(), "http: 404 Not Found")
}

func TestAsApiError(t *testing.T) {
	client := tCreateClient(t)

	req, err := client.NewRequest(context.Background(), "GET", "models", nil)

	tAssertNil(t, err)

	req.URL.RawQuery = "context=0"

	_, err = client.Do(req)

	tAssertNotNil(t, err)
	tAssertEquals(t, err.Error(), "api error ZodError: Too small: expected number to be >=1")
}

func TestAsProviderError(t *testing.T) {
	client := tCreateClient(t)

	data := ChatCompletionRequest{
		Model: "openai/gpt-5.6-luna",
		Messages: []ChatMessage{
			{
				Role: ChatRoleUser,
				Content: ChatContent{
					Text: "hi",
				},
			},
		},
		ResponseFormat: &ChatResponseFormat{
			Type: ChatResponseFormatTypeJSONObject,
		},
	}

	_, err := client.CreateChatCompletion(context.Background(), data)

	tAssertNotNil(t, err)
	tAssertEquals(t, err.Error(), "provider error: Response input messages must contain the word 'json' in some form to use 'text.format' of type 'json_object'.")
}
