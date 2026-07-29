package openingrouter

import (
	"context"
	"testing"
)

func TestListEmbeddingModels(t *testing.T) {
	client := tCreateClient(t)

	models, err := client.ListEmbeddingModels(context.Background(), &ListEmbeddingModelsOptions{
		Offset: new(0),
		Limit:  new(10),
	})

	tAssertNil(t, err)
	tAssertLen(t, models, 10)

	model := models[0]

	tAssertMinLen(t, model.ID, 3)
	tAssertMinLen(t, model.Name, 3)
}
