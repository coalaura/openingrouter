package openingrouter

import (
	"context"
	"testing"
)

func TestGetModelBySlug(t *testing.T) {
	client := tCreateClient(t)

	model, err := client.GetModelBySlug(context.Background(), ExampleTextModel)

	tAssertNil(t, err)
	tAssertMinLen(t, model.ID, 3)
	tAssertMinLen(t, model.Name, 3)
}

func TestListModels(t *testing.T) {
	client := tCreateClient(t)

	models, err := client.ListModels(context.Background(), &ListModelsOptions{
		Offset: new(0),
		Limit:  new(10),
	})

	tAssertNil(t, err)
	tAssertLen(t, models, 10)

	model := models[0]

	tAssertMinLen(t, model.ID, 3)
	tAssertMinLen(t, model.Name, 3)
}

func TestListUserModels(t *testing.T) {
	client := tCreateClient(t)

	models, err := client.ListUserModels(context.Background(), &ListUserModelsOptions{
		Offset: new(0),
		Limit:  new(10),
	})

	tAssertNil(t, err)
	tAssertLen(t, models, 10)

	model := models[0]

	tAssertMinLen(t, model.ID, 3)
	tAssertMinLen(t, model.Name, 3)
}
