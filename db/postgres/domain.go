package postgres

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/models"
)

// ReadDomain reads an undeleted domain from the database by uuid
func (c *Client) ReadDomain(id uuid.UUID) (*models.Domain, error) {
	var domain models.Domain
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
func (c *Client) ReadDomainZ(id uuid.UUID) (*models.Domain, error) {
	var domain models.Domain
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
func (c *Client) ReadDomainByDomain(d string) (*models.Domain, error) {
	var domain models.Domain
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
func (c *Client) ReadDomainsForUser(userID uuid.UUID) (*[]models.Domain, error) {
	var domains []models.Domain
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
