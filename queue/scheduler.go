package queue

import "github.com/google/uuid"

type Scheduler interface {
	AddDomain(uuid.UUID) error
	RemoveDomain(uuid.UUID) error
}
