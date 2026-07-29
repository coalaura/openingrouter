package openingrouter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
)

type streamErrorCarrier interface {
	streamError() error
}

// OpenrouterStream represents a stream of typed elements.
type OpenrouterStream[T any] interface {
	Recv() (T, error)
	Close()
}

// ServerSentEventsStream receives Server-Sent Events from an HTTP response.
type ServerSentEventsStream[T any] struct {
	stream    <-chan T
	done      chan struct{}
	response  *http.Response
	closeOnce sync.Once
}

// JsonResponseStream wraps a slice of chunks as a stream.
type JsonResponseStream[T any] struct {
	chunks []T
	index  atomic.Uint64
	closed atomic.Bool
}

// Recv reads the next chunk from the stream.
func (s *ServerSentEventsStream[T]) Recv() (T, error) {
	select {
	case chunk, ok := <-s.stream:
		if ok {
			carrier, ok := any(chunk).(streamErrorCarrier)
			if ok {
				err := carrier.streamError()
				if err != nil {
					var zero T

					return zero, err
				}
			}

			return chunk, nil
		}
	case <-s.done:
	}

	var zero T

	return zero, io.EOF
}

// Close terminates the stream and cleans up resources.
func (s *ServerSentEventsStream[T]) Close() {
	s.closeOnce.Do(func() {
		close(s.done)

		if s.response != nil {
			s.response.Body.Close()
		}
	})
}

// Recv returns the next chunk, or io.EOF when exhausted or after Close.
func (s *JsonResponseStream[T]) Recv() (T, error) {
	var zero T

	if s.closed.Load() {
		return zero, io.EOF
	}

	i := s.index.Add(1) - 1
	if i >= uint64(len(s.chunks)) {
		return zero, io.EOF
	}

	return s.chunks[i], nil
}

// Close makes subsequent Recv calls return io.EOF.
func (s *JsonResponseStream[T]) Close() {
	s.closed.Store(true)
}

// Add appends one or more chunks to the stream.
func (s *JsonResponseStream[T]) Add(chunks ...T) {
	s.chunks = append(s.chunks, chunks...)
}

// NewServerSentEventsStream starts a background reader over resp.Body and
// yields decoded T values. The returned stream owns resp: call Close (typically
// via defer) to cancel the reader and release the body.
func NewServerSentEventsStream[T any](ctx context.Context, resp *http.Response) *ServerSentEventsStream[T] {
	out := make(chan T)
	done := make(chan struct{})

	sse := &ServerSentEventsStream[T]{
		stream:   out,
		done:     done,
		response: resp,
	}

	go func() {
		defer close(out)
		defer sse.Close()

		reader := bufio.NewReader(resp.Body)
		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			default:
			}

			line, err := reader.ReadBytes('\n')
			if err != nil {
				return
			}

			trimmed := bytes.TrimSpace(line)
			if len(trimmed) == 0 || trimmed[0] == ':' {
				continue
			}

			payload := bytes.TrimSpace(bytes.TrimPrefix(trimmed, []byte("data:")))
			if len(payload) == 0 || bytes.Equal(payload, []byte("[DONE]")) {
				return
			}

			var chunk T

			err = json.Unmarshal(payload, &chunk)
			if err != nil {
				return
			}

			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			case out <- chunk:
			}
		}
	}()

	return sse
}

// NewJsonResponseStream returns a new JsonResponseStream initialized with chunks.
func NewJsonResponseStream[T any](chunks ...T) *JsonResponseStream[T] {
	return &JsonResponseStream[T]{
		chunks: chunks,
	}
}

// IsResponseServerSentEventsStream reports whether the HTTP response is a Server-Sent Events stream.
func IsResponseServerSentEventsStream(response *http.Response) bool {
	contentType := strings.ToLower(response.Header.Get("Content-Type"))

	return strings.Contains(contentType, "text/event-stream")
}
