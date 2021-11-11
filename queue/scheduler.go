package queue

import "github.com/google/uuid"

// Scheduler adds jobs to the queue
type Scheduler interface {
	AddDomain(uuid.UUID) error
	RemoveDomain(uuid.UUID) error
}
