package openingrouter

import (
	"encoding/json"
	"io"
	"strings"
	"testing"
)

type mediaMessageTestCase struct {
	name     string
	message  ChatMessage
	expected string
}

func TestMediaMessageRequest(t *testing.T) {
	agenticMessage := UserMessageWithVideo("Describe this", "https://example.com/video.mp4")
	agenticMessage.Content.Parts[1].VideoURL.Processing = ChatVideoProcessingAgentic

	staticMessage := UserMessageWithVideo("Describe this", "data:video/webm;base64,GkXfow==")
	staticMessage.Content.Parts[1].VideoURL.Processing = ChatVideoProcessingStatic

	tests := []mediaMessageTestCase{
		{
			name:     "wav audio",
			message:  UserMessageWithAudio("Transcribe this", "UklGRg==", "wav"),
			expected: `{"messages":[{"role":"user","content":[{"type":"text","text":"Transcribe this"},{"type":"input_audio","input_audio":{"data":"UklGRg==","format":"wav"}}]}]}`,
		},
		{
			name:     "mp3 audio",
			message:  UserMessageWithAudio("Transcribe this", "SUQz", "mp3"),
			expected: `{"messages":[{"role":"user","content":[{"type":"text","text":"Transcribe this"},{"type":"input_audio","input_audio":{"data":"SUQz","format":"mp3"}}]}]}`,
		},
		{
			name:     "video URL",
			message:  UserMessageWithVideo("Describe this", "https://example.com/video.mp4"),
			expected: `{"messages":[{"role":"user","content":[{"type":"text","text":"Describe this"},{"type":"video_url","video_url":{"url":"https://example.com/video.mp4"}}]}]}`,
		},
		{
			name:     "video data URL",
			message:  UserMessageWithVideo("Describe this", "data:video/webm;base64,GkXfow=="),
			expected: `{"messages":[{"role":"user","content":[{"type":"text","text":"Describe this"},{"type":"video_url","video_url":{"url":"data:video/webm;base64,GkXfow=="}}]}]}`,
		},
		{
			name:     "agentic video processing",
			message:  agenticMessage,
			expected: `{"messages":[{"role":"user","content":[{"type":"text","text":"Describe this"},{"type":"video_url","video_url":{"url":"https://example.com/video.mp4","processing":"agentic"}}]}]}`,
		},
		{
			name:     "static video processing",
			message:  staticMessage,
			expected: `{"messages":[{"role":"user","content":[{"type":"text","text":"Describe this"},{"type":"video_url","video_url":{"url":"data:video/webm;base64,GkXfow==","processing":"static"}}]}]}`,
		},
	}

	client := NewClient("test-token")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := ChatCompletionRequest{
				Messages: []ChatMessage{test.message},
			}

			req, err := client.NewRequest(t.Context(), "POST", "chat/completions", request)
			tAssertNil(t, err)

			defer req.Body.Close()

			body, err := io.ReadAll(req.Body)
			tAssertNil(t, err)
			tAssertEquals(t, strings.TrimSpace(string(body)), test.expected)

			var decoded ChatCompletionRequest

			err = json.Unmarshal(body, &decoded)
			tAssertNil(t, err)
			tAssertEquals(t, len(decoded.Messages), 1)
			tAssertEquals(t, decoded.Messages[0].Content.String(), test.message.Content.String())

			encoded, err := json.Marshal(decoded)
			tAssertNil(t, err)
			tAssertEquals(t, string(encoded), test.expected)
		})
	}
}
