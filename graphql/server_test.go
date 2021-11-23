package graphql

import (
	"net/http"
	"testing"
)

type mockResponseWriter struct {
	Body   []byte
	Headr  http.Header
	Status int
}

func (w *mockResponseWriter) Header() http.Header {
	return w.Headr
}

func (w *mockResponseWriter) Write(b []byte) (int, error) {
	w.Body = append(w.Body, b...)
	return len(b), nil
}

func (w *mockResponseWriter) WriteHeader(status int) {
	w.Status = status
	return
}

func TestNewServer(t *testing.T) {
	server, _, _, _, err := newTestServer()

	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
		return
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
		return
	}
}
