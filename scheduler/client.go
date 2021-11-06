package scheduler

import (
	faktory "github.com/contribsys/faktory/client"
)

// Client is Faktory client used to queue jobs for workers.
type Client struct {
	faktory *faktory.Client
}

// NewClient creates a new Faktory client
func NewClient() (*Client, error) {
	client := Client{}
	var err error

	client.faktory, err = faktory.Open()
	if err != nil {
		return nil, err
	}
	return &client, nil
}
