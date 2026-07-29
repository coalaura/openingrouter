package openingrouter

import (
	"context"
	"testing"
)

func TestCreateEmbeddings(t *testing.T) {
	client := tCreateClient(t)

	data := EmbeddingRequest{
		Model: "openai/text-embedding-3-small",
		Input: EmbeddingInput{
			Text: "Hello World",
		},
	}

	resp, err := client.CreateEmbeddings(context.Background(), data)

	tAssertNil(t, err)
	tAssertEquals(t, resp.Object, EmbeddingObjectList)
	tAssertLen(t, resp.Data, 1)

	embedding := resp.Data[0]

	tAssertEquals(t, embedding.Object, EmbeddingObjectEmbedding)
	tAssertEquals(t, embedding.Index, 0)
	tAssertMinLen(t, embedding.Embedding.Floats, 1)

	usage := resp.Usage

	tAssertNotNil(t, usage)
	tAssertNotNil(t, usage.Cost)

	testUsage += *usage.Cost
}
