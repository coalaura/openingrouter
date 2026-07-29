package openingrouter

// STTRequest represents the request body of the speech-to-text transcription
// endpoint. Model and InputAudio are required, zero values and nil pointers of
// the remaining fields are omitted from the request.
type STTRequest struct {
	Model      string        `json:"model"`
	InputAudio STTInputAudio `json:"input_audio"`

	Language               string                    `json:"language,omitempty"`
	Provider               *STTProviderPreferences   `json:"provider,omitempty"`
	ResponseFormat         STTResponseFormat         `json:"response_format,omitempty"`
	Temperature            *float64                  `json:"temperature,omitempty"`
	TimestampGranularities []STTTimestampGranularity `json:"timestamp_granularities,omitempty"`
}

// STTInputAudio holds base64-encoded audio data and its format to transcribe.
type STTInputAudio struct {
	Data   string `json:"data"`
	Format string `json:"format"`
}

// STTProviderPreferences represents the provider specific passthrough
// configuration of a speech-to-text request.
type STTProviderPreferences struct {
	Options ProviderOptions `json:"options,omitempty"`
}

// STTResponse is the root response returned by the speech-to-text transcription
// endpoint. Text is required and always present. Additional fields such as
// Duration, Language, Segments, Task, and Words are populated when
// ResponseFormat is set to verbose_json.
type STTResponse struct {
	Text string `json:"text"`

	Duration *float64     `json:"duration,omitempty"`
	Language string       `json:"language,omitempty"`
	Segments []STTSegment `json:"segments,omitempty"`
	Task     string       `json:"task,omitempty"`
	Usage    *STTUsage    `json:"usage,omitempty"`
	Words    []STTWord    `json:"words,omitempty"`
}

// STTSegment represents a timestamped transcript segment, returned when
// ResponseFormat is verbose_json.
type STTSegment struct {
	ID               int     `json:"id"`
	Seek             int     `json:"seek"`
	Start            float64 `json:"start"`
	End              float64 `json:"end"`
	Text             string  `json:"text"`
	Tokens           []int   `json:"tokens,omitempty"`
	Temperature      float64 `json:"temperature"`
	AvgLogprob       float64 `json:"avg_logprob"`
	CompressionRatio float64 `json:"compression_ratio"`
	NoSpeechProb     float64 `json:"no_speech_prob"`
	Speaker          *int    `json:"speaker,omitempty"`
}

// STTWord represents a timestamped word, returned when the provider includes
// word-level timestamps.
type STTWord struct {
	Word    string  `json:"word"`
	Start   float64 `json:"start"`
	End     float64 `json:"end"`
	Speaker *int    `json:"speaker,omitempty"`
}

// STTUsage represents aggregated usage statistics for a speech-to-text request.
type STTUsage struct {
	Cost         *float64 `json:"cost,omitempty"`
	InputTokens  int      `json:"input_tokens,omitempty"`
	OutputTokens int      `json:"output_tokens,omitempty"`
	Seconds      *float64 `json:"seconds,omitempty"`
	TotalTokens  int      `json:"total_tokens,omitempty"`
}

// STTResponseFormat is the output response format of a transcription request.
type STTResponseFormat string

const (
	STTResponseFormatJSON        STTResponseFormat = "json"
	STTResponseFormatVerboseJSON STTResponseFormat = "verbose_json"
)

// STTTimestampGranularity is a timestamp detail level for verbose_json
// transcription responses.
type STTTimestampGranularity string

const (
	STTTimestampGranularityWord    STTTimestampGranularity = "word"
	STTTimestampGranularitySegment STTTimestampGranularity = "segment"
)
