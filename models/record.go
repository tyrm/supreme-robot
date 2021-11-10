package models

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

const (
	// RecordTypeA is the type for an A type record
	RecordTypeA = "A"
	// RecordTypeAAAA is the type for an A type record
	RecordTypeAAAA = "AAAA"
)

// Record is a dns record.
type Record struct {
	Name     string    `db:"name"`
	DomainID uuid.UUID `db:"domain_id"`
	Type     string    `db:"type"`
	Value    string    `db:"value"`
	TTL      int       `db:"ttl"`

	Priority sql.NullInt32  `db:"priority"`
	Port     sql.NullInt32  `db:"port"`
	Weight   sql.NullInt32  `db:"weight"`
	Refresh  sql.NullInt32  `db:"refresh"`
	Retry    sql.NullInt32  `db:"retry"`
	Expire   sql.NullInt32  `db:"expire"`
	MBox     sql.NullString `db:"mbox"`
	Tag      sql.NullString `db:"tag"`

	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Validate validates the data provided
func (r *Record) Validate() error {
	switch r.Type {
	case RecordTypeA:
		// check for required attributes
		if r.Name == "" {
			return errMissingName
		}
		if r.Value == "" {
			return errMissingIP
		}
		if r.TTL == 0 {
			return errMissingTTL
		}

		// check values
		if !reSubDomain.MatchString(r.Name) {
			return errInvalidName
		}
		if !reIPv4Address.MatchString(r.Value) {
			return errInvalidIP
		}
		if r.TTL < 1 {
			return errInvalidTTL
		}

		return nil
	case RecordTypeAAAA:
		// check for required attributes
		if r.Name == "" {
			return errMissingName
		}
		if r.Value == "" {
			return errMissingIP
		}
		if r.TTL == 0 {
			return errMissingTTL
		}

		// check values
		if !reSubDomain.MatchString(r.Name) {
			return errInvalidName
		}
		if !reIPv6Address.MatchString(r.Value) {
			return errInvalidIP
		}
		if r.TTL < 1 {
			return errInvalidTTL
		}

		return nil
	case "":
		return errMissingType
	default:
		return errUnknownType
	}
}
