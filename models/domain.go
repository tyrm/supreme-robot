package models

import (
	"database/sql"
	"github.com/google/uuid"
	"regexp"
	"strings"
	"time"
)

var reValidDomain = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-.]*\\.$")

// Domain represents a domain name to be served
type Domain struct {
	Domain  string    `db:"domain" json:"domain"`
	OwnerID uuid.UUID `db:"owner_id" json:"-"`

	Records *[]Record `db:"-" json:"records"`

	ID        uuid.UUID `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// Model Functions

func (d *Domain) create(c *Client) error {
	// add to database
	err := c.db.
		QueryRowx(`INSERT INTO "public"."domains"("domain", "owner_id")
			VALUES ($1, $2) RETURNING id, created_at, updated_at;`, d.Domain, d.OwnerID).
		Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)

	return err
}

func (d *Domain) delete(c *Client) error {
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

// GetRecords retrieves the Records for the domain from the database
func (d *Domain) GetRecords(c *Client) (*[]Record, error) {
	records, err := c.ReadRecordsForDomain(d, "name", true)
	if err != nil {
		return nil, err
	}
	return records, nil
}

// ValidateDomain checks that the domain name doesn't contain any invalid values
func (d *Domain) ValidateDomain() bool {
	return reValidDomain.MatchString(d.Domain)
}

// Client Functions

// ReadDomain reads an undeleted domain from the database by uuid
func (c *Client) ReadDomain(id uuid.UUID) (*Domain, error) {
	var domain Domain
	err := c.db.
		Get(&domain, `SELECT id, domain, owner_id, created_at, updated_at
		FROM public.domains WHERE id = $1 AND deleted_at IS NULL;`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &domain, nil
}

// ReadDomainZ reads an any domain from the database by uuid even after deleted by user
func (c *Client) ReadDomainZ(id uuid.UUID) (*Domain, error) {
	var domain Domain
	err := c.db.
		Get(&domain, `SELECT id, domain, owner_id, created_at, updated_at
		FROM public.domains WHERE id = $1;`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &domain, nil
}

// ReadDomainByDomain will read a domain from the database by domain name.
func (c *Client) ReadDomainByDomain(d string) (*Domain, error) {
	var domain Domain
	err := c.db.
		Get(&domain, `SELECT id, domain, owner_id, created_at, updated_at
		FROM public.domains WHERE lower(domain) = lower($1) AND deleted_at IS NULL;`, d)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &domain, nil
}

// ReadDomainsForUser will read all domains for a given user.
func (c *Client) ReadDomainsForUser(userID uuid.UUID) (*[]Domain, error) {
	var domains []Domain
	err := c.db.
		Select(&domains, `SELECT id, domain, owner_id, created_at, updated_at
		FROM public.domains WHERE owner_id = $1 AND deleted_at IS NULL;`, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &domains, nil
}

// ReadDomainsPageForUser will read a sorted and paginated list of domains for a given user.
func (c *Client) ReadDomainsPageForUser(user *User, index, count int, orderBy string, asc bool) (*[]Domain, error) {
	var domainList []Domain

	// build query
	query := "SELECT id, domain, owner_id, created_at, updated_at FROM public.domains WHERE owner_id = $1 AND deleted_at IS NULL ORDER BY "

	switch strings.ToLower(orderBy) {
	case "created_at":
		query = query + "created_at "
	case "domain":
		query = query + "domain "
	default:
		return nil, errUnknownAttribute
	}

	if asc {
		query = query + "ASC "
	} else {
		query = query + "DESC "
	}

	query = query + "OFFSET $2 LIMIT $3;"

	// run query
	offset := index * count
	err := c.db.Select(&domainList, query, user.ID, offset, count)
	if err != nil && err != sql.ErrNoRows {
		logger.Errorf("cant get user page: %s")
		return nil, err
	}

	return &domainList, nil
}
