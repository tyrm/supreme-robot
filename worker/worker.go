package worker

import (
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/kv"
	"github.com/tyrm/supreme-robot/queue"
)

// Worker is a
type Worker struct {
	db      db.DB
	manager queue.Manager
	kv      kv.DNS
}

// Close closes the worker
func (w *Worker) Close() error {
	w.manager.Terminate(false)
	return nil
}

// Run runs the worker
func (w *Worker) Run() error {
	w.manager.Run()
	return nil
}

// NewWorker creates a new faktory worker
func NewWorker(k kv.DNS, m queue.Manager, d db.DB) (*Worker, error) {
	worker := Worker{
		db:      d,
		manager: m,
		kv:      k,
	}

	worker.manager.ProcessStrictPriorityQueues(queue.QueueDNS)

	worker.manager.Register(queue.JobAddDomain, worker.addDomainHandler)
	worker.manager.Register(queue.JobRemoveDomain, worker.removeDomainHandler)
	worker.manager.Register(queue.JobUpdateSubDomain, worker.updateSubDomainHandler)

	return &worker, nil
}
