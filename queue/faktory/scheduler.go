package faktory

import (
	faktory "github.com/contribsys/faktory/client"
)

// Scheduler is Faktory client used to queue jobs for workers.
type Scheduler struct {
	faktory *faktory.Client
}

// NewClient creates a new Faktory client
func NewClient() (*Scheduler, error) {
	client := Scheduler{}
	var err error

	client.faktory, err = faktory.Open()
	if err != nil {
		return nil, err
	}
	return &client, nil
}
