package faktory

import (
	faktory "github.com/contribsys/faktory/client"
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/queue"
)

// AddDomain adds a job to add a new domain to redis.
func (c *Client) AddDomain(id uuid.UUID) error {
	job := faktory.NewJob(queue.JobAddDomain, id)
	job.Queue = queue.QueueDNS
	return c.faktory.Push(job)
}

// RemoveDomain adds a job to purge a domain from redis.
func (c *Client) RemoveDomain(id uuid.UUID) error {
	job := faktory.NewJob(queue.JobRemoveDomain, id)
	job.Queue = queue.QueueDNS
	return c.faktory.Push(job)
}

// UpdateSubDomain adds a job to updates redis with records from db for a given subdomain
func (c *Client) UpdateSubDomain(id uuid.UUID, name string) error {
	job := faktory.NewJob(queue.JobUpdateSubDomain, id, name)
	job.Queue = queue.QueueDNS
	return c.faktory.Push(job)
}
