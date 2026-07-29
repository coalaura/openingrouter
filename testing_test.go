package openingrouter

import (
	"os"
	"reflect"
	"testing"
)

const (
	ExampleTextModel = "deepseek/deepseek-v4-flash"
)

type iterable interface {
	~string | ~[]FrontendModel | ~[]Model | ~[]ImageModel | ~[]GeneratedImage | ~[]ChatChoice
}

func tCreateClient(t testing.TB) *Client {
	t.Helper()

	token := os.Getenv("OPENROUTER_API_KEY")
	if token == "" {
		t.Skip("OPENROUTER_API_KEY not set")
	}

	return NewClient(token)
}

func tAssertLen[T iterable](t testing.TB, value T, length int) {
	t.Helper()

	if len(value) == length {
		return
	}

	t.Fatalf("expected len(%T) to be %d, got: %d", value, length, len(value))
}

func tAssertMinLen[T iterable](t testing.TB, value T, length int) {
	t.Helper()

	if len(value) >= length {
		return
	}

	t.Fatalf("expected len(%T) to be at least %d, got: %d", value, length, len(value))
}

func tAssertEquals[T comparable](t testing.TB, actual, expected T) {
	t.Helper()

	if expected == actual {
		return
	}

	t.Fatalf("expected %v, got: %v", expected, actual)
}

func tAssertNil(t testing.TB, actual any) {
	t.Helper()

	if actual == nil {
		return
	}

	val := reflect.ValueOf(actual)

	switch val.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		if val.IsNil() {
			return
		}
	}

	t.Fatalf("expected nil, got: %#v", actual)
}

func tAssertNotNil(t testing.TB, actual any) {
	t.Helper()

	if actual != nil {
		val := reflect.ValueOf(actual)

		switch val.Kind() {
		case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
			if !val.IsNil() {
				return
			}
		default:
			return
		}
	}

	t.Fatal("expected not nil, got: nil")
}
