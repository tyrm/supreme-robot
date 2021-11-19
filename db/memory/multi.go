package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"time"
)

// CreateDomainWRecords will create a domain and it's records in a single database transaction.
func (c *Client) CreateDomainWRecords(domain *models.Domain, records ...*models.Record) error {
	c.Lock()
	defer c.Unlock()

	now := time.Now()

	domain.ID = uuid.New()
	domain.CreatedAt = now
	domain.UpdatedAt = now

	recordList := make([]models.Record, len(records))
	for i, r := range records {
		r.DomainID = domain.ID
		r.ID = uuid.New()
		r.CreatedAt = now
		r.UpdatedAt = now

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
	u, err := c.ReadUser(userID)
	if err != nil {
		return err
	}
	if u == nil {
		return db.ErrUnknownUser
	}

	c.Lock()
	defer c.Unlock()

	for _, group := range groupIDs {
		u.Groups = append(u.Groups, group)
	}

	c.users[userID] = *u

	return nil
}
