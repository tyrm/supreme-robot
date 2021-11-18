package memory

import (
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"time"
)

// Update a struct
func (c *Client) Update(obj interface{}) error {
	switch obj := obj.(type) {
	case *models.User:
		return c.updateUser(obj)
	default:
		return db.ErrUnknownType
	}
}

func (c *Client) updateUser(u *models.User) error {
	u.UpdatedAt = time.Now()

	// Lock DB
	c.Lock()
	defer c.Unlock()

	c.users[u.ID] = *u
	return nil
}
