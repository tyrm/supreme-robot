package postgres

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/models"
)

// ReadUser will retrieve a user by their uuid from the database
func (c *Client) ReadUser(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := c.db.
		Get(&user, `SELECT id ,username, password, created_at, updated_at 
		FROM public.users WHERE id = $1 AND deleted_at IS NULL;`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var groups []uuid.UUID
	err = c.db.
		Select(&groups, `SELECT group_id 
		FROM public.group_membership WHERE user_id = $1 AND deleted_at IS NULL;`, user.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	user.Groups = groups

	return &user, nil
}

// ReadUserByUsername will read a user by username from the database
func (c *Client) ReadUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := c.db.
		Get(&user, `SELECT id ,username, password, created_at, updated_at 
		FROM public.users WHERE lower(username) = lower($1) AND deleted_at IS NULL;`, username)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var groups []uuid.UUID
	err = c.db.
		Select(&groups, `SELECT group_id 
		FROM public.group_membership WHERE user_id = $1 AND deleted_at IS NULL;`, user.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	user.Groups = groups

	return &user, nil
}
