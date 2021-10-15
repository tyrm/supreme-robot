package models

import (
	"database/sql"
	"time"
)

type Domain struct {
	Domain string `db:"domain" json:"domain"`
	OwnerID string `db:"owner_id" json:"-"`

	Owner *User `db:"-" json:"owner"`

	ID        string    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Model Functions
func (d *Domain) Create(c *Client) error {
	var err error

	// add to database
	if d.ID == "" {
		// id doesn't exist
		err = c.db.
			QueryRowx(`INSERT INTO "public"."domains"("domain", "owner_id")
			VALUES ($1, $2) RETURNING id, created_at, updated_at;`, d.Domain, d.Owner.ID).
			Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
	} else {
		// id exists
		err = c.db.
			QueryRowx(`INSERT INTO "public"."domains"("id", "domain", "owner_id")
			VALUES ($1, $2, $3) RETURNING created_at, updated_at;`, d.ID, d.Domain, d.Owner.ID).
			Scan(&d.CreatedAt, &d.UpdatedAt)
	}

	return err
}

// Client Functions
func (c *Client) ReadDomainsForUser(user *User) (*[]Domain, error) {
	var domains []Domain
	err := c.db.
		Select(&domains, `SELECT id, domain, owner_id, created_at, updated_at
		FROM public.domains WHERE owner_id = $1 AND deleted_at IS NULL;`, user.ID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &domains, nil
}