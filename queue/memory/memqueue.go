package memory

import (
	faktory "github.com/contribsys/faktory_worker_go"
	"github.com/juju/loggo"
	"github.com/tyrm/supreme-robot/queue"
)

var logger = loggo.GetLogger("scheduler.mem")

// MemQueue is a database client.
type MemQueue struct {
	Callbacks     map[string]faktory.Perform
	QueuesEnabled map[string]bool
	Queues        map[string]chan []interface{}
}

// NewMemQueue creates a new models Client from Config
func NewMemQueue() (*MemQueue, error) {
	logger.Tracef("starting memory scheduler")

	s := MemQueue{
		Callbacks:     make(map[string]faktory.Perform),
		QueuesEnabled: make(map[string]bool),
		Queues:        make(map[string]chan []interface{}),
	}

	s.Queues[queue.QueueDNS] = make(chan []interface{}, 1024)

	return &s, nil
}
