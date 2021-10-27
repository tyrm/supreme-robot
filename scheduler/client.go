package scheduler

import (
	faktory "github.com/contribsys/faktory/client"
)

type Client struct {
	faktory *faktory.Client
}

func NewClient() (*Client, error) {
	client := Client{}
	var err error

	client.faktory, err = faktory.Open()
	if err != nil {
		return nil, err
	}
	return &client, nil
}
