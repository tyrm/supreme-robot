package memory

import (
	"github.com/juju/loggo"
)

var logger = loggo.GetLogger("scheduler.mem")

// Scheduler is a database client.
type Scheduler struct {
}

// NewScheduler creates a new models Client from Config
func NewScheduler() (*Scheduler, error) {
	logger.Tracef("starting memory scheduler")

	s := Scheduler{}

	return &s, nil
}
