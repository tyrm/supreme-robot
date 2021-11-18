package db

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/models"
)

// DB represents the required commands for a database.
type DB interface {
	Create(interface{}) error
	CreateDomainWRecords(*models.Domain, ...*models.Record) error
	CreateGroupsForUser(uuid.UUID, ...uuid.UUID) error
	Delete(interface{}) error
	ReadDomain(uuid.UUID) (*models.Domain, error)
	ReadDomainZ(uuid.UUID) (*models.Domain, error)
	ReadDomainByDomain(string) (*models.Domain, error)
	ReadDomainsForUser(uuid.UUID) (*[]models.Domain, error)
	ReadRecordsForDomain(uuid.UUID, string, bool) (*[]models.Record, error)
	ReadUser(uuid.UUID) (*models.User, error)
	ReadUserByUsername(string) (*models.User, error)
	Update(interface{}) error
}
