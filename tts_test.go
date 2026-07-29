package openingrouter

import (
	"context"
	"io"
	"testing"
)

func TestCreateSpeech(t *testing.T) {
	client := tCreateClient(t)

	data := SpeechRequest{
		Model:          "sesame/csm-1b",
		Input:          "hello world",
		ResponseFormat: SpeechResponseFormatMP3,
	}

	resp, err := client.CreateSpeech(context.Background(), data)

	tAssertNil(t, err)

	defer resp.Body.Close()

	tAssertEquals(t, resp.ContentType, "audio/mpeg")

	var header [3]byte

	n, err := io.ReadFull(resp.Body, header[:])

	tAssertNil(t, err)

	tAssertMP3Header(t, header[:n])
}
