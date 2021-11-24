package memory

import (
	"github.com/juju/loggo"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var logger = loggo.GetLogger("db.mem")

// Client is a database client.
type Client struct {
	sync.RWMutex

	KV *cache.Cache
}

// NewClient creates a new models Client from Config
func NewClient() (*Client, error) {
	logger.Tracef("starting memory kv")

	c := Client{
		KV: cache.New(cache.NoExpiration, 10*time.Minute),
	}

	return &c, nil
}
