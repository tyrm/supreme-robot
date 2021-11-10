package graphql

import (
	"github.com/tyrm/supreme-robot/config"
	dbMem "github.com/tyrm/supreme-robot/db/memory"
	kvMem "github.com/tyrm/supreme-robot/kv/memory"
	queueMem "github.com/tyrm/supreme-robot/queue/memory"
	"net/http"
	"testing"
	"time"
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
	ws, _, _, _, err := newTestServer()

	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if ws == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}
}

func newTestServer() (*Server, *queueMem.Scheduler, *dbMem.Client, *kvMem.Client, error) {
	cnf := config.Config{
		AccessExpiration:  time.Hour * 24,
		AccessSecret:      "test",
		PrimaryNS:         "ns1.example.com.",
		RefreshExpiration: time.Hour * 24,
		RefreshSecret:     "test1234",
	}

	db, err := dbMem.NewClient()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	kv, err := kvMem.NewClient()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	qc, err := queueMem.NewScheduler()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	ws, err := NewServer(&cnf, qc, db, kv)

	return ws, qc, db, kv, nil
}
