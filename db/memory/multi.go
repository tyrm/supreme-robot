package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/models"
)

// CreateDomainWRecords will create a domain and it's records in a single database transaction.
func (c *Client) CreateDomainWRecords(domain *models.Domain, records ...*models.Record) error {

	return nil
}

// CreateGroupsForUser adds group_membership entries for the user to the database
func (c *Client) CreateGroupsForUser(userID uuid.UUID, groupIDs ...uuid.UUID) error {

	return nil
}
