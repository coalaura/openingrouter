package openingrouter

import (
	"context"
	"testing"
)

func TestGetCurrentApiKey(t *testing.T) {
	client := tCreateClient(t)

	info, err := client.GetCurrentApiKey(context.Background())

	tAssertNil(t, err)
	tAssertNotNil(t, info)

	tAssertMinLen(t, info.Label, 1)
}
