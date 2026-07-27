package openingrouter

import (
	"context"
	"errors"
	"io"
	"testing"
)

func TestGenerateImage(t *testing.T) {
	client := tCreateClient(t)

	data := ImageGenerationRequest{
		Model:       "black-forest-labs/flux.2-klein-4b",
		Prompt:      "Generate an image of a cat in a banana costume.",
		AspectRatio: ImageAspectRatio1x1,
	}

	resp, err := client.GenerateImage(context.Background(), data)

	tAssertNil(t, err)
	tAssertLen(t, resp.Data, 1)

	image := resp.Data[0]

	tAssertMinLen(t, image.MediaType, 8)
	tAssertMinLen(t, image.B64JSON, 32)
}

func TestGenerateImageStream(t *testing.T) {
	client := tCreateClient(t)

	data := ImageGenerationRequest{
		Model:       "black-forest-labs/flux.2-klein-4b",
		Prompt:      "Generate an image of a cat in a banana costume.",
		AspectRatio: ImageAspectRatio1x1,
		Stream:      new(true),
	}

	stream, err := client.GenerateImageStream(context.Background(), data)

	tAssertNil(t, err)

	var result *ImageStreamEvent

	for {
		entry, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if entry.Type == ImageStreamEventTypeCompleted {
			result = &entry
		}
	}

	tAssertNotNil(t, result)
	tAssertMinLen(t, result.MediaType, 8)
	tAssertMinLen(t, result.B64JSON, 32)
}
