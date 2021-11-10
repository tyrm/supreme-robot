package worker

import (
	faktory "github.com/contribsys/faktory_worker_go"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/redis"
	"github.com/tyrm/supreme-robot/scheduler"
)

// Worker is a
type Worker struct {
	db      db.DB
	manager *faktory.Manager
	redis   *redis.Client
}

// Run runs the worker
func (w *Worker) Run() error {
	w.manager.Run()
	return nil
}

// NewWorker creates a new faktory worker
func NewWorker(r *redis.Client, d db.DB) (*Worker, error) {
	worker := Worker{
		db:      d,
		manager: faktory.NewManager(),
		redis:   r,
	}

	worker.manager.ProcessStrictPriorityQueues("default", scheduler.QueueDNS)

	worker.manager.Register("AddDomain", worker.addDomainHandler)
	worker.manager.Register("RemoveDomain", worker.removeDomainHandler)

	return &worker, nil
}
