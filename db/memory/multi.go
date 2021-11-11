package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/models"
)

// CreateDomainWRecords will create a domain and it's records in a single database transaction.
func (c *Client) CreateDomainWRecords(domain *models.Domain, records ...*models.Record) error {
	c.Lock()
	defer c.Unlock()

	domain.ID = uuid.New()

	recordList := make([]models.Record, len(records))
	for i, r := range records {
		r.DomainID = domain.ID
		r.ID = uuid.New()

		c.records[r.ID] = *r
		recordList[i] = *r
	}
	domain.Records = &recordList

	c.domains[domain.ID] = *domain
	c.domainsZ[domain.ID] = *domain
	return nil
}

// CreateGroupsForUser adds group_membership entries for the user to the database
func (c *Client) CreateGroupsForUser(userID uuid.UUID, groupIDs ...uuid.UUID) error {

	return nil
}
