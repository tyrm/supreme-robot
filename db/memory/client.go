package memory

import (
	"database/sql"
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

//
var now = time.Now()
var userAdmin = models.User{
	ID:       uuid.MustParse("44892097-2c97-4c16-b4d1-e8522586df48"),
	Username: "admin",
	Password: "$2a$14$mmOFu7eOyQUFC0S/gopbDeJKcADiUx7QleU85WW7FnnCiXNgENb1G", // password
	Groups: []uuid.UUID{
		models.GroupSuperAdmin,
	},
	CreatedAt: now,
	UpdatedAt: now,
}
var userUser = models.User{
	ID:        uuid.MustParse("69706e10-fe5a-4bd3-a494-75d5c4230f5f"),
	Username:  "user",
	Password:  "$2a$14$eqjfDrVLAhpa9tivYSRQSOLnOGmrwjXQURlC79CZpDT/WY3HdqGRO", // starship
	Groups:    []uuid.UUID{},
	CreatedAt: now,
	UpdatedAt: now,
}

var domainTest1 = models.Domain{
	ID:        uuid.MustParse("b5c54740-807d-433c-8be9-73adcd794352"),
	Domain:    "test1.",
	OwnerID:   userAdmin.ID,
	CreatedAt: now,
	UpdatedAt: now,
}
var domainTest2 = models.Domain{
	ID:        uuid.MustParse("fa84d136-90b9-4fb4-a9fb-fb73a4a918aa"),
	Domain:    "test2.",
	OwnerID:   userAdmin.ID,
	CreatedAt: now,
	UpdatedAt: now,
}
var domainTest3 = models.Domain{
	ID:        uuid.MustParse("ab6fb9d0-b87e-4c76-87c8-4ec07d3fb3c2"),
	Domain:    "test3.",
	OwnerID:   userAdmin.ID,
	CreatedAt: now,
	UpdatedAt: now,
}
var domainExample = models.Domain{
	ID:        uuid.MustParse("eee6cac1-031d-4fe6-8a2b-c4ef2dd79695"),
	Domain:    "example.",
	OwnerID:   userUser.ID,
	CreatedAt: now,
	UpdatedAt: now,
}

var recordTest1SOA = models.Record{
	ID:       uuid.MustParse("d546c0ab-5d8b-4d7c-9bf9-5b4e5a8947a3"),
	Name:     "@",
	DomainID: domainTest1.ID,
	Type:     models.RecordTypeSOA,
	Value:    "ns1.example.com.",
	TTL:      300,
	MBox: sql.NullString{
		String: "hostmaster.test1.",
		Valid:  true,
	},
	Refresh: sql.NullInt32{
		Int32: 22,
		Valid: true,
	},
	Retry: sql.NullInt32{
		Int32: 44,
		Valid: true,
	},
	Expire: sql.NullInt32{
		Int32: 33,
		Valid: true,
	},
	CreatedAt: now,
	UpdatedAt: now,
}
var recordTest1NS1 = models.Record{
	ID:        uuid.MustParse("42cd351e-1287-4f9d-8c55-4acdd9623951"),
	Name:      "@",
	DomainID:  domainTest1.ID,
	Type:      models.RecordTypeNS,
	Value:     "ns1.example.com.",
	TTL:       300,
	CreatedAt: now,
	UpdatedAt: now,
}
var recordTest2SOA = models.Record{
	ID:       uuid.MustParse("a2887d9d-2133-4270-9544-b716e2156857"),
	Name:     "@",
	DomainID: domainTest2.ID,
	Type:     models.RecordTypeSOA,
	Value:    "ns1.example.com.",
	TTL:      300,
	MBox: sql.NullString{
		String: "hostmaster.test2.",
		Valid:  true,
	},
	Refresh: sql.NullInt32{
		Int32: 22,
		Valid: true,
	},
	Retry: sql.NullInt32{
		Int32: 44,
		Valid: true,
	},
	Expire: sql.NullInt32{
		Int32: 33,
		Valid: true,
	},
	CreatedAt: now,
	UpdatedAt: now,
}
var recordTest2NS1 = models.Record{
	ID:        uuid.MustParse("d0884a6f-af31-42ae-b342-b17ffe9dbdb7"),
	Name:      "@",
	DomainID:  domainTest2.ID,
	Type:      models.RecordTypeNS,
	Value:     "ns1.example.com.",
	TTL:       300,
	CreatedAt: now,
	UpdatedAt: now,
}
var recordTest3SOA = models.Record{
	ID:       uuid.MustParse("5fe7d098-e0fa-477c-b954-4f265343a83b"),
	Name:     "@",
	DomainID: domainTest3.ID,
	Type:     models.RecordTypeSOA,
	Value:    "ns1.example.com.",
	TTL:      300,
	MBox: sql.NullString{
		String: "hostmaster.test3.",
		Valid:  true,
	},
	Refresh: sql.NullInt32{
		Int32: 22,
		Valid: true,
	},
	Retry: sql.NullInt32{
		Int32: 44,
		Valid: true,
	},
	Expire: sql.NullInt32{
		Int32: 33,
		Valid: true,
	},
	CreatedAt: now,
	UpdatedAt: now,
}
var recordTest3NS1 = models.Record{
	ID:        uuid.MustParse("edfba0fd-6c62-41c6-96cd-bacd03958c72"),
	Name:      "@",
	DomainID:  domainTest3.ID,
	Type:      models.RecordTypeNS,
	Value:     "ns1.example.com.",
	TTL:       300,
	CreatedAt: now,
	UpdatedAt: now,
}
var recordExampleSOA = models.Record{
	ID:       uuid.MustParse("fcbadd89-33ab-43c6-bfb2-10a421b94225"),
	Name:     "@",
	DomainID: domainExample.ID,
	Type:     models.RecordTypeSOA,
	Value:    "ns1.example.com.",
	TTL:      300,
	MBox: sql.NullString{
		String: "hostmaster.example.com.",
		Valid:  true,
	},
	Refresh: sql.NullInt32{
		Int32: 22,
		Valid: true,
	},
	Retry: sql.NullInt32{
		Int32: 44,
		Valid: true,
	},
	Expire: sql.NullInt32{
		Int32: 33,
		Valid: true,
	},
	CreatedAt: now,
	UpdatedAt: now,
}
var recordExampleNS1 = models.Record{
	ID:        uuid.MustParse("4bf21a38-7144-48b7-ac74-8cba9b3d5a65"),
	Name:      "@",
	DomainID:  domainExample.ID,
	Type:      models.RecordTypeNS,
	Value:     "ns1.example.com.",
	TTL:       300,
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
	c.users[userUser.ID] = userUser

	// add test domains
	c.domains[domainTest1.ID] = domainTest1
	c.domainsZ[domainTest1.ID] = domainTest1
	c.domains[domainTest2.ID] = domainTest2
	c.domainsZ[domainTest2.ID] = domainTest2
	c.domains[domainTest3.ID] = domainTest3
	c.domainsZ[domainTest3.ID] = domainTest3
	c.domains[domainExample.ID] = domainExample
	c.domainsZ[domainExample.ID] = domainExample

	// add records
	c.records[recordTest1SOA.ID] = recordTest1SOA
	c.records[recordTest1NS1.ID] = recordTest1NS1
	c.records[recordTest2SOA.ID] = recordTest2SOA
	c.records[recordTest2NS1.ID] = recordTest2NS1
	c.records[recordTest3SOA.ID] = recordTest3SOA
	c.records[recordTest3NS1.ID] = recordTest3NS1
	c.records[recordExampleSOA.ID] = recordExampleSOA
	c.records[recordExampleNS1.ID] = recordExampleNS1

	return &c, nil
}
