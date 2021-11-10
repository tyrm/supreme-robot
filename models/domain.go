package models

import (
	"github.com/google/uuid"
	"time"
)

// Domain represents a domain name to be served
type Domain struct {
	Domain  string    `db:"domain" json:"domain"`
	OwnerID uuid.UUID `db:"owner_id" json:"-"`

	Records *[]Record `db:"-" json:"records"`

	ID        uuid.UUID `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// Validate checks that the domain name doesn't contain any invalid values
func (d *Domain) Validate() bool {
	return reTopDomain.MatchString(d.Domain)
}
