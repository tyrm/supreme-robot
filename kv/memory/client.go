package memory

import (
	"github.com/juju/loggo"
)

var logger = loggo.GetLogger("db.mem")

// Client is a database client.
type Client struct {
	KV map[string]string
}

// NewClient creates a new models Client from Config
func NewClient() (*Client, error) {
	logger.Tracef("starting memory kv")

	c := Client{
		KV: make(map[string]string),
	}

	return &c, nil
}
