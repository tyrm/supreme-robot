package graphql

import (
	"github.com/tyrm/supreme-robot/config"
	dbMem "github.com/tyrm/supreme-robot/db/memory"
	kvMem "github.com/tyrm/supreme-robot/kv/memory"
	queueMem "github.com/tyrm/supreme-robot/queue/memory"
	"time"
)

var (
	testServer      *Server
	testServerQueue *queueMem.MemQueue
	testServerDB    *dbMem.Client
	testServerKv    *kvMem.Client
)

func newTestServer() (*Server, error) {
	cnf := config.Config{
		AccessExpiration:  time.Hour * 24,
		AccessSecret:      "test",
		HTTPPort:          ":27413",
		PrimaryNS:         "ns1.example.com.",
		RefreshExpiration: time.Hour * 24,
		RefreshSecret:     "test1234",
	}

	if testServerQueue == nil {
		qc, err := queueMem.NewMemQueue()
		if err != nil {
			return nil, err
		}
		testServerQueue = qc
	}

	if testServerDB == nil {
		db, err := dbMem.NewClient()
		if err != nil {
			return nil, err
		}
		testServerDB = db
	}

	if testServerKv == nil {
		kv, err := kvMem.NewClient()
		if err != nil {
			return nil, err
		}
		testServerKv = kv
	}

	if testServer == nil {
		ws, err := NewServer(&cnf, testServerQueue, testServerDB, testServerKv)
		if err != nil {
			return nil, err
		}
		testServer = ws
	}

	return testServer, nil
}
