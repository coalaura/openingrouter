package openingrouter

import (
	"context"
	"testing"
)

func TestListFrontendModels(t *testing.T) {
	list, err := ListFrontendModels(context.Background())

	tAssertNil(t, err)
	tAssertMinLen(t, list, 100)
}
