package memory

import (
	"github.com/google/uuid"
	"github.com/juju/loggo"
	"github.com/tyrm/supreme-robot/models"
	"sync"
	"time"
)

var logger = loggo.GetLogger("db.mem")

// Client is a database client.
type Client struct {
	sync.RWMutex

	domains  map[uuid.UUID]models.Domain
	domainsZ map[uuid.UUID]models.Domain
	records  map[uuid.UUID]models.Record
	users    map[uuid.UUID]models.User
}

var now = time.Now()
var userAdmin = models.User{
	ID:       uuid.MustParse("8c504483-1e11-4243-b6c8-14499877a641"),
	Username: "admin",
	Password: "$2a$14$mmOFu7eOyQUFC0S/gopbDeJKcADiUx7QleU85WW7FnnCiXNgENb1G", // password
	Groups: []uuid.UUID{
		models.GroupSuperAdmin,
	},
	CreatedAt: now,
	UpdatedAt: now,
}

// NewClient creates a new models Client from Config
func NewClient() (*Client, error) {
	logger.Tracef("starting memory db")

	c := Client{
		domains:  make(map[uuid.UUID]models.Domain),
		domainsZ: make(map[uuid.UUID]models.Domain),
		records:  make(map[uuid.UUID]models.Record),
		users:    make(map[uuid.UUID]models.User),
	}

	// add test users
	c.users[userAdmin.ID] = userAdmin

	return &c, nil
}
