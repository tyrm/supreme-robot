package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"strings"
)

// ReadRecordsForDomain will read all records for a given domain
func (c *Client) ReadRecordsForDomain(domainID uuid.UUID, orderBy string, asc bool) (*[]models.Record, error) {
	switch strings.ToLower(orderBy) {
	case "created_at":
	case "name":
	default:
		return nil, db.ErrUnknownAttribute
	}

	foundRecords := make([]models.Record, 0)

	// Lock DB
	c.RLock()
	defer c.RUnlock()

	for _, record := range c.records {
		if record.DomainID == domainID {
			foundRecords = append(foundRecords, record)
		}
	}

	return &foundRecords, nil
}
