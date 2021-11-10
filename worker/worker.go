package worker

import (
	faktory "github.com/contribsys/faktory_worker_go"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/kv"
	faktory2 "github.com/tyrm/supreme-robot/scheduler"
)

// Worker is a
type Worker struct {
	db      db.DB
	manager *faktory.Manager
	kv      kv.DNS
}

// Run runs the worker
func (w *Worker) Run() error {
	w.manager.Run()
	return nil
}

// NewWorker creates a new faktory worker
func NewWorker(k kv.DNS, d db.DB) (*Worker, error) {
	worker := Worker{
		db:      d,
		manager: faktory.NewManager(),
		kv:      k,
	}

	worker.manager.ProcessStrictPriorityQueues("default", faktory2.QueueDNS)

	worker.manager.Register("AddDomain", worker.addDomainHandler)
	worker.manager.Register("RemoveDomain", worker.removeDomainHandler)

	return &worker, nil
}
