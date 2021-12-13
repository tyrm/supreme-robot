package queue

import (
	faktory "github.com/contribsys/faktory_worker_go"
)

// Manager issues jobs based on the queue
type Manager interface {
	ProcessStrictPriorityQueues(...string)
	Register(string, faktory.Perform)
	Run()
	Terminate(bool)
}
