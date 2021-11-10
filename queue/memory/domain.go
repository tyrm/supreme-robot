package memory

import (
	"github.com/google/uuid"
)

// AddDomain adds a job to add a new domain to redis.
func (s *Scheduler) AddDomain(id uuid.UUID) error {
	return nil
}

// RemoveDomain adds a job to purge a domain from redis.
func (s *Scheduler) RemoveDomain(id uuid.UUID) error {
	return nil
}
