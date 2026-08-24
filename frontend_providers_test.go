package openingrouter

import (
	"context"
	"testing"
)

func TestListFrontendProviders(t *testing.T) {
	list, err := ListFrontendProviders(context.Background())

	tAssertNil(t, err)
	tAssertMinLen(t, list, 50)
}
