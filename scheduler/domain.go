package scheduler

import (
	faktory "github.com/contribsys/faktory/client"
	"github.com/google/uuid"
)

func (c *Client) AddDomain(id uuid.UUID) error {
	job := faktory.NewJob("AddDomain", id)
	return c.faktory.Push(job)
}

func (c *Client) RemoveDomain(id uuid.UUID) error {
	job := faktory.NewJob("RemoveDomain", id)
	return c.faktory.Push(job)
}
