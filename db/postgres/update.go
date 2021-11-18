package postgres

import (
	"database/sql"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
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
	err := c.db.
		QueryRowx(`UPDATE "public"."users"
		SET password=$1, updated_at=CURRENT_TIMESTAMP
		WHERE id=$2;`, u.Password, u.ID).
		Scan()
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}
