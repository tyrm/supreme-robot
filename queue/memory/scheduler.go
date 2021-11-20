package memory

import (
	"github.com/juju/loggo"
	"sync"
)

var logger = loggo.GetLogger("scheduler.mem")

// Scheduler is a database client.
type Scheduler struct {
	sync.RWMutex

	Jobs map[string][][]interface{}
}

// NewScheduler creates a new models Client from Config
func NewScheduler() (*Scheduler, error) {
	logger.Tracef("starting memory scheduler")

	s := Scheduler{
		Jobs: make(map[string][][]interface{}),
	}

	return &s, nil
}
