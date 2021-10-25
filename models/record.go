package models

import (
	"database/sql"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Record struct {
	Name     string    `db:"name"`
	DomainID uuid.UUID `db:"domain_id"`
	Type     string    `db:"type"`
	Value    string    `db:"value"`

	TTL      sql.NullInt32  `db:"ttl"`
	Priority sql.NullInt32  `db:"priority"`
	Port     sql.NullInt32  `db:"port"`
	Weight   sql.NullInt32  `db:"weight"`
	Refresh  sql.NullInt32  `db:"refresh"`
	Retry    sql.NullInt32  `db:"retry"`
	Expire   sql.NullInt32  `db:"expire"`
	MBox     sql.NullString `db:"mbox"`
	Tag      sql.NullString `db:"tag"`

	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Model Functions
func (r *Record) Create(c *Client) error {
	var err error

	// add to database
	if r.ID == uuid.Nil {
		// id doesn't exist
		err = c.db.
			QueryRowx(`INSERT INTO "public"."domain_records"("name", "domain_id", "type", "value", "ttl",
            "priority", "port", "weight", "refresh", "retry", "expire", "mbox", "tag")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id, created_at, updated_at;`,
			r.Name, r.DomainID, r.Type, r.Value, r. TTL, r.Priority, r.Port, r.Weight, r.Refresh, r.Retry, r.Expire,
			r.MBox, r.Tag).
			Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt)
	} else {
		// id exists
		err = c.db.
			QueryRowx(`INSERT INTO "public"."domain_records"("id", "name", "domain_id", "type", "value", "ttl",
            "priority", "port", "weight", "refresh", "retry", "expire", "mbox", "tag")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id, created_at, updated_at;`,
			r.ID, r.Name, r.DomainID, r.Type, r.Value, r. TTL, r.Priority, r.Port, r.Weight, r.Refresh, r.Retry,
			r.Expire, r.MBox, r.Tag).
			Scan(&r.CreatedAt, &r.UpdatedAt)
	}

	return err
}
func (r *Record) Domain(c *Client) (*Domain, error) {
	return c.ReadDomain(r.DomainID)
}

// Client Functions
func (c *Client) ReadRecord(id uuid.UUID) (*Record, error) {
	var record Record
	err := c.db.
		Get(&record, `SELECT id, name, domain_id, type, value, ttl, priority, port, weight, refresh, retry,
        expire, mbox, tag, created_at, updated_at 
		FROM public.domain_records WHERE id = $1 AND deleted_at IS NULL;`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &record, nil
}


func (c *Client) ReadRecordsPageForDomain(domain *Domain, index, count int, orderBy string, asc bool) (*[]Record, error) {
	var recordList []Record

	// build query
	query := "SELECT id, name, domain_id, type, value, ttl, priority, port, weight, refresh, retry, " +
		"expire, mbox, tag, created_at, updated_at  FROM public.domains " +
		"WHERE domain_id = $1 AND deleted_at IS NULL ORDER BY "

	switch strings.ToLower(orderBy) {
	case "created_at":
		query = query + "created_at "
	case "name":
		query = query + "name "
	default:
		return nil, ErrUnknownAttribute
	}

	if asc {
		query = query + "ASC "
	} else {
		query = query + "DESC "
	}

	query = query + "OFFSET $2 LIMIT $3;"

	// run query
	offset := index * count
	err := c.db.Select(&recordList, query, domain.ID, offset, count)
	if err != nil && err != sql.ErrNoRows {
		logger.Errorf("cant get record page: %s")
		return nil, err
	}

	return &recordList, nil
}