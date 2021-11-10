package faktory

import (
	faktory "github.com/contribsys/faktory/client"
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/scheduler"
)

// AddDomain adds a job to add a new domain to redis.
func (c *Client) AddDomain(id uuid.UUID) error {
	job := faktory.NewJob("AddDomain", id)
	job.Queue = scheduler.QueueDNS
	return c.faktory.Push(job)
}

// RemoveDomain adds a job to purge a domain from redis.
func (c *Client) RemoveDomain(id uuid.UUID) error {
	job := faktory.NewJob("RemoveDomain", id)
	job.Queue = scheduler.QueueDNS
	return c.faktory.Push(job)
}
