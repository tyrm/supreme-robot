package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/models"
	"strings"
)

// ReadUser will retrieve a user by their uuid from the database
func (c *Client) ReadUser(id uuid.UUID) (*models.User, error) {
	// Lock DB
	c.RLock()
	defer c.RUnlock()

	user, userOk := c.users[id]
	if !userOk {
		return nil, nil
	}
	return &user, nil
}

// ReadUserByUsername will read a user by username from the database
func (c *Client) ReadUserByUsername(username string) (*models.User, error) {
	// Lock DB
	c.RLock()
	defer c.RUnlock()

	for _, v := range c.users {
		if strings.ToLower(username) == strings.ToLower(v.Username) {
			return &v, nil
		}
	}

	return nil, nil
}
