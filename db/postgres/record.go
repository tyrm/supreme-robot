package postgres

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"strings"
)

// ReadRecordsForDomain will read all records for a given domain
func (c *Client) ReadRecordsForDomain(domainID uuid.UUID, orderBy string, asc bool) (*[]models.Record, error) {
	var recordList []models.Record

	// build query
	query := "SELECT id, name, domain_id, type, value, ttl, priority, port, weight, refresh, retry, " +
		"expire, mbox, tag, created_at, updated_at FROM public.domain_records " +
		"WHERE domain_id = $1 AND deleted_at IS NULL ORDER BY "

	switch strings.ToLower(orderBy) {
	case "created_at":
		query = query + "created_at "
	case "name":
		query = query + "name "
	default:
		return nil, db.ErrUnknownAttribute
	}

	if asc {
		query = query + "ASC, id ASC;"
	} else {
		query = query + "DESC, id DESC;"
	}

	// run query
	err := c.db.Select(&recordList, query, domainID)
	if err != nil && err != sql.ErrNoRows {
		logger.Errorf("cant get record page: %s", err.Error())
		return nil, err
	}

	return &recordList, nil
}

// ReadRecordsForDomainByName will read records for a given domain matching the given name
func (c *Client) ReadRecordsForDomainByName(domainID uuid.UUID, name string) (*[]models.Record, error) {
	var recordList []models.Record

	// build query
	query := "SELECT id, name, domain_id, type, value, ttl, priority, port, weight, refresh, retry, " +
		"expire, mbox, tag, created_at, updated_at FROM public.domain_records " +
		"WHERE domain_id = $1 AND lower(name) = lower($2) AND deleted_at IS NULL ORDER BY created_at ASC, id ASC;"

	// run query
	err := c.db.Select(&recordList, query, domainID, name)
	if err != nil && err != sql.ErrNoRows {
		logger.Errorf("cant get record page: %s", err.Error())
		return nil, err
	}

	return &recordList, nil
}
