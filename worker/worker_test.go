package worker

import (
	dbMem "github.com/tyrm/supreme-robot/db/memory"
	kvMem "github.com/tyrm/supreme-robot/kv/memory"
	queueMem "github.com/tyrm/supreme-robot/queue/memory"
	"testing"
)

var (
	testWorker      *Worker
	testWorkerQueue *queueMem.MemQueue
	testWorkerDB    *dbMem.Client
	testWorkerKv    *kvMem.Client
)

func TestNewWorker(t *testing.T) {
	server, err := newTestWorker()

	if err != nil {
		t.Errorf("unexpected error, got: %s, want: error.", err.Error())
		return
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Worker.")
		return
	}
}

func newTestWorker() (*Worker, error) {
	if testWorkerQueue == nil {
		qc, err := queueMem.NewMemQueue()
		if err != nil {
			return nil, err
		}
		testWorkerQueue = qc
	}

	if testWorkerDB == nil {
		db, err := dbMem.NewClient()
		if err != nil {
			return nil, err
		}
		testWorkerDB = db
	}

	if testWorkerKv == nil {
		kv, err := kvMem.NewClient()
		if err != nil {
			return nil, err
		}
		testWorkerKv = kv
	}

	if testWorker == nil {
		ws, err := NewWorker(testWorkerKv, testWorkerQueue, testWorkerDB)
		if err != nil {
			return nil, err
		}
		testWorker = ws
	}

	return testWorker, nil
}
