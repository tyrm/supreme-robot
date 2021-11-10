package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/models"
)

// ReadDomain reads an undeleted domain from the database by uuid
func (c *Client) ReadDomain(id uuid.UUID) (*models.Domain, error) {
	// Lock DB
	c.RLock()
	defer c.RUnlock()

	return nil, nil
}

// ReadDomainZ reads an any domain from the database by uuid even after deleted by user
func (c *Client) ReadDomainZ(id uuid.UUID) (*models.Domain, error) {
	// Lock DB
	c.RLock()
	defer c.RUnlock()

	return nil, nil
}

// ReadDomainByDomain will read a domain from the database by domain name.
func (c *Client) ReadDomainByDomain(d string) (*models.Domain, error) {
	// Lock DB
	c.RLock()
	defer c.RUnlock()

	return nil, nil
}

// ReadDomainsForUser will read all domains for a given user.
func (c *Client) ReadDomainsForUser(userID uuid.UUID) (*[]models.Domain, error) {
	// Lock DB
	c.RLock()
	defer c.RUnlock()

	return nil, nil
}
