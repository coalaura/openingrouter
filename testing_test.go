package openingrouter

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

const (
	ExampleTextModel = "deepseek/deepseek-v4-flash"
)

type iterable interface {
	~string | ~[]byte | ~[]float64 |
		~[]FrontendModel | ~[]Model | ~[]ImageModel |
		~[]GeneratedImage | ~[]ChatChoice | ~[]Embedding
}

var testUsage float64

func TestMain(m *testing.M) {
	code := m.Run()

	fmt.Fprintf(os.Stderr, "test usage: $%.6f\n", testUsage)

	os.Exit(code)
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

func tAssertContainsFold(t testing.TB, actual, contains string) {
	t.Helper()

	if strings.Contains(strings.ToLower(actual), contains) {
		return
	}

	t.Fatalf("expected %q to contain %q", actual, contains)
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

func tAssertMP3Header(t testing.TB, header []byte) {
	t.Helper()

	if len(header) >= 3 {
		// ID3v2 container header
		if header[0] == 'I' && header[1] == 'D' && header[2] == '3' {
			return
		}

		// Raw MPEG frame sync
		if header[0] == 0xFF && (header[1]&0xE0) == 0xE0 {
			return
		}
	}

	t.Fatalf("expected mp3 header, got: %v", header)
}
