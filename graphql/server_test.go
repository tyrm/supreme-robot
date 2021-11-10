package graphql

import (
	"github.com/tyrm/supreme-robot/config"
	dbMem "github.com/tyrm/supreme-robot/db/memory"
	kvMem "github.com/tyrm/supreme-robot/kv/memory"
	queueMem "github.com/tyrm/supreme-robot/queue/memory"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	cnf := config.Config{
		AccessExpiration:  time.Hour * 24,
		AccessSecret:      "test",
		PrimaryNS:         "ns1.example.com.",
		RefreshExpiration: time.Hour * 24,
		RefreshSecret:     "test1234",
	}

	db, err := dbMem.NewClient()
	if err != nil {
		t.Errorf("expected error, got: nil, want: error.")
		return
	}

	kv, err := kvMem.NewClient()
	if err != nil {
		t.Errorf("expected error, got: nil, want: error.")
		return
	}

	qc, err := queueMem.NewScheduler()
	if err != nil {
		t.Errorf("expected error, got: nil, want: error.")
		return
	}

	// create web server
	ws, err := NewServer(&cnf, qc, db, kv)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
	}
	if ws == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}
}
