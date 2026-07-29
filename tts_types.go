package openingrouter

import "io"

// SpeechRequest represents the request body of the speech synthesis endpoint.
// Model, Input and Voice are required, zero values and nil pointers of the
// remaining fields are omitted from the request. Voice identifiers are provider
// specific, the supported set of a model is listed in Model.SupportedVoices.
type SpeechRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
	Voice string `json:"voice"`

	Provider       *SpeechProviderPreferences `json:"provider,omitempty"`
	ResponseFormat SpeechResponseFormat       `json:"response_format,omitempty"`
	Speed          *float64                   `json:"speed,omitempty"`
}

// SpeechProviderPreferences represents the provider specific passthrough
// configuration of a speech request. Unlike ProviderPreferences the speech
// endpoint accepts passthrough options only, there is no provider routing.
type SpeechProviderPreferences struct {
	Options ProviderOptions `json:"options,omitempty"`
}

// SpeechResponse represents the synthesized audio of a speech request. Body is
// the raw audio bytestream and is owned by the caller, it must be closed once
// read. ContentType is the media type of the bytestream and varies by the
// requested format (audio/mpeg for mp3, audio/pcm for 16-bit little-endian pcm).
type SpeechResponse struct {
	ContentType string
	Body        io.ReadCloser
}

// SpeechResponseFormat is the audio encoding of a synthesized bytestream. It
// defaults to pcm.
type SpeechResponseFormat string

const (
	SpeechResponseFormatMP3 SpeechResponseFormat = "mp3"
	SpeechResponseFormatPCM SpeechResponseFormat = "pcm"
)
