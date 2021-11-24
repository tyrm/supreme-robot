package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/queue"
)

// AddDomain adds a job to add a new domain to redis.
func (s *Scheduler) AddDomain(id uuid.UUID) error {
	newJob := []interface{}{
		queue.JobAddDomain,
		id.String(),
	}

	s.Lock()
	defer s.Unlock()

	s.Jobs[queue.QueueDNS] = append(s.Jobs[queue.QueueDNS], newJob)
	return nil
}

// RemoveDomain adds a job to purge a domain from redis.
func (s *Scheduler) RemoveDomain(id uuid.UUID) error {
	newJob := []interface{}{
		queue.JobRemoveDomain,
		id.String(),
	}

	s.Lock()
	defer s.Unlock()

	s.Jobs[queue.QueueDNS] = append(s.Jobs[queue.QueueDNS], newJob)
	return nil
}

// UpdateSubDomain adds a job to updates redis with records from db for a given subdomain
func (s *Scheduler) UpdateSubDomain(id uuid.UUID, name string) error {
	newJob := []interface{}{
		queue.JobUpdateSubDomain,
		id.String(),
		name,
	}

	s.Lock()
	defer s.Unlock()

	s.Jobs[queue.QueueDNS] = append(s.Jobs[queue.QueueDNS], newJob)
	return nil
}
