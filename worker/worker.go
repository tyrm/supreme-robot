package worker

import (
	faktory "github.com/contribsys/faktory_worker_go"
	"github.com/tyrm/supreme-robot/models"
	"github.com/tyrm/supreme-robot/redis"
)

type Worker struct {
	db      *models.Client
	manager *faktory.Manager
	redis   *redis.Client
}

func (w *Worker) Run() error {
	w.manager.Run()
	return nil
}

func NewWorker(r *redis.Client, d *models.Client) (*Worker, error) {
	worker := Worker{
		db:      d,
		manager: faktory.NewManager(),
		redis:   r,
	}

	worker.manager.Register("AddDomain", worker.addDomainHandler)
	worker.manager.Register("RemoveDomain", worker.removeDomainHandler)

	return &worker, nil
}
