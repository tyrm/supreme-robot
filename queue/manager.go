package queue

import (
	faktory "github.com/contribsys/faktory_worker_go"
)

type Manager interface {
	ProcessStrictPriorityQueues(...string)
	Register(string, faktory.Perform)
	Run()
}
