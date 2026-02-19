package openingrouter

import (
	"testing"
)

func tAssertNoErr(t testing.TB, err error) {
	t.Helper()

	if err == nil {
		return
	}

	t.Fatalf("tAssertNoErr failed, got: %v", err)
}

func tAssertMinLen[T any](t testing.TB, sl []T, ln int) {
	t.Helper()

	if len(sl) >= ln {
		return
	}

	t.Fatalf("tAssertMinLen failed, expected at least %d, got: %d", ln, len(sl))
}
