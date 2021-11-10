package memory

import (
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
)

// Delete a struct
func (c *Client) Delete(obj interface{}) error {
	switch obj := obj.(type) {
	case *models.Domain:
		return c.deleteDomain(obj)
	default:
		return db.ErrUnknownType
	}
}

func (c *Client) deleteDomain(d *models.Domain) error {
	return nil
}
