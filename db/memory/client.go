package memory

import (
	"github.com/google/uuid"
	"github.com/juju/loggo"
	"github.com/tyrm/supreme-robot/models"
	"sync"
)

var logger = loggo.GetLogger("db.mem")

// Client is a database client.
type Client struct {
	sync.RWMutex

	domains map[uuid.UUID]models.Domain
	records map[uuid.UUID]models.Record
	users   map[uuid.UUID]models.User
}

// NewClient creates a new models Client from Config
func NewClient() (*Client, error) {
	logger.Tracef("starting memory db")

	c := Client{
		domains: make(map[uuid.UUID]models.Domain),
		records: make(map[uuid.UUID]models.Record),
		users:   make(map[uuid.UUID]models.User),
	}

	adminUser := models.User{
		ID:       uuid.Must(uuid.Parse("44892097-2c97-4c16-b4d1-e8522586df48")),
		Username: "admin",
	}
	err := adminUser.SetPassword("password")
	if err != nil {
		return nil, err
	}
	c.users[adminUser.ID] = adminUser

	return &c, nil
}
