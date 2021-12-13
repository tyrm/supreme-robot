package memory

import (
	faktoryClient "github.com/contribsys/faktory/client"
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/queue"
)

// AddDomain adds a job to add a new domain to redis.
func (m *MemQueue) AddDomain(id uuid.UUID) error {
	newJob := []interface{}{
		queue.JobAddDomain,
		faktoryClient.RandomJid(),
		id.String(),
	}

	m.Queues[queue.QueueDNS] <- newJob
	return nil
}

// RemoveDomain adds a job to purge a domain from redis.
func (m *MemQueue) RemoveDomain(id uuid.UUID) error {
	newJob := []interface{}{
		queue.JobRemoveDomain,
		faktoryClient.RandomJid(),
		id.String(),
	}

	m.Queues[queue.QueueDNS] <- newJob
	return nil
}

// UpdateSubDomain adds a job to updates redis with records from db for a given subdomain
func (m *MemQueue) UpdateSubDomain(id uuid.UUID, name string) error {
	newJob := []interface{}{
		queue.JobUpdateSubDomain,
		faktoryClient.RandomJid(),
		id.String(),
		name,
	}

	m.Queues[queue.QueueDNS] <- newJob
	return nil
}
