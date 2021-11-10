package postgres

import (
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
	// add to database
	return c.db.
		QueryRowx(`INSERT INTO "public"."domains"("domain", "owner_id")
			VALUES ($1, $2) RETURNING id, created_at, updated_at;`, d.Domain, d.OwnerID).
		Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
}

func (c *Client) createRecord(r *models.Record) error {
	// add to database
	return c.db.
		QueryRowx(`INSERT INTO "public"."domain_records"("name", "domain_id", "type", "value", "ttl",
            "priority", "port", "weight", "refresh", "retry", "expire", "mbox", "tag")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id, created_at, updated_at;`,
			r.Name, r.DomainID, r.Type, r.Value, r.TTL, r.Priority, r.Port, r.Weight, r.Refresh, r.Retry, r.Expire,
			r.MBox, r.Tag).
		Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt)
}

func (c *Client) createUser(u *models.User) error {
	// add to database
	return c.db.
		QueryRowx(`INSERT INTO "public"."users"("username", "password")
			VALUES ($1, $2) RETURNING id, created_at, updated_at;`, u.Username, u.Password).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}
