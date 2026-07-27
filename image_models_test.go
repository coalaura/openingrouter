package openingrouter

import (
	"context"
	"testing"
)

func TestListImageModels(t *testing.T) {
	client := tCreateClient(t)

	models, err := client.ListImageModels(context.Background())

	tAssertNil(t, err)
	tAssertMinLen(t, models, 1)

	model := models[0]

	tAssertMinLen(t, model.ID, 3)
	tAssertMinLen(t, model.Name, 3)
	tAssertMinLen(t, model.Endpoints, 3)
}
