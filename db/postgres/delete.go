package postgres

import (
	"database/sql"
	"github.com/tyrm/supreme-robot/models"
)

// Delete a struct
func (c *Client) Delete(obj interface{}) error {
	switch obj := obj.(type) {
	case *models.Domain:
		return c.deleteDomain(obj)
	default:
		return errUnknownType
	}
}

func (c *Client) deleteDomain(d *models.Domain) error {
	err := c.db.
		QueryRowx(`UPDATE "public"."domains"
		SET deleted_at=CURRENT_TIMESTAMP
		WHERE id=$1;`, d.ID).
		Scan()
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}
