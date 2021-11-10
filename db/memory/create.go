package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
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
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	d.ID = id

	return nil
}

func (c *Client) createRecord(r *models.Record) error {
	return nil
}

func (c *Client) createUser(u *models.User) error {
	return nil
}
