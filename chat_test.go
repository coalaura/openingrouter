package openingrouter

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestCreateChatCompletion(t *testing.T) {
	client := tCreateClient(t)

	data := ChatCompletionRequest{
		Model: "openai/gpt-oss-20b",
		Messages: []ChatMessage{
			{
				Role: ChatRoleUser,
				Content: ChatContent{
					Text: "Respond exactly with \"hello world\". No quotes, no newlines or whitespace, do not output anything else.",
				},
			},
		},
	}

	resp, err := client.CreateChatCompletion(context.Background(), data)

	tAssertNil(t, err)
	tAssertLen(t, resp.Choices, 1)

	message := resp.Choices[0].Message

	tAssertEquals(t, message.Role, ChatRoleAssistant)
	tAssertEquals(t, message.Content.Text, "hello world")
}

func TestCreateChatCompletionStream(t *testing.T) {
	client := tCreateClient(t)

	data := ChatCompletionRequest{
		Model: "openai/gpt-oss-20b",
		Messages: []ChatMessage{
			{
				Role: ChatRoleUser,
				Content: ChatContent{
					Text: "Respond exactly with \"hello world\". No quotes, no newlines or whitespace, do not output anything else.",
				},
			},
		},
		Stream: new(true),
	}

	stream, err := client.CreateChatCompletionStream(context.Background(), data)

	tAssertNil(t, err)

	var (
		role   ChatRole
		result strings.Builder
	)

	for {
		entry, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		choices := entry.Choices

		if len(choices) > 0 {
			delta := choices[0].Delta

			if role == "" && delta.Role != "" {
				role = delta.Role
			}

			if delta.Content != "" {
				result.WriteString(delta.Content)
			}
		}
	}

	tAssertEquals(t, role, ChatRoleAssistant)
	tAssertEquals(t, result.String(), "hello world")
}
