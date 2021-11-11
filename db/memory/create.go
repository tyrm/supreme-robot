package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"time"
)

// Create a struct
func (c *Client) Create(obj interface{}) error {
	switch obj := obj.(type) {
	case *models.Domain:
		return c.createDomain(obj)
	case *models.Record:
		return c.createRecord(obj)
	case *models.User:
		return c.createUser(obj)
	default:
		return db.ErrUnknownType
	}
}

func (c *Client) createDomain(d *models.Domain) error {
	now := time.Now()

	d.ID = uuid.New()
	d.CreatedAt = now
	d.UpdatedAt = now

	// Lock DB
	c.Lock()
	defer c.Unlock()

	c.domains[d.ID] = *d
	c.domainsZ[d.ID] = *d

	return nil
}

func (c *Client) createRecord(r *models.Record) error {
	now := time.Now()

	r.ID = uuid.New()
	r.CreatedAt = now
	r.UpdatedAt = now

	// Lock DB
	c.Lock()
	defer c.Unlock()

	c.records[r.ID] = *r
	return nil
}

func (c *Client) createUser(u *models.User) error {
	now := time.Now()

	u.ID = uuid.New()
	u.CreatedAt = now
	u.UpdatedAt = now

	// Lock DB
	c.Lock()
	defer c.Unlock()

	c.users[u.ID] = *u
	return nil
}
