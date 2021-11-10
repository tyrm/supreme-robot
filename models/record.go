package models

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

const (
	// RecordTypeA is the type for an A type record
	RecordTypeA = "A"
	// RecordTypeAAAA is the type for an AAAA type record
	RecordTypeAAAA = "AAAA"
	// RecordTypeCNAME is the type for a CNAME type record
	RecordTypeCNAME = "CNAME"
	// RecordTypeMX is the type for an MX type record
	RecordTypeMX = "MX"
	// RecordTypeNS is the type for a NS type record
	RecordTypeNS = "NS"
	// RecordTypeTXT is the type for a TXT type record
	RecordTypeTXT = "TXT"
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
	case RecordTypeCNAME:
		// check for required attributes
		if r.Name == "" {
			return errMissingName
		}
		if r.Value == "" {
			return errMissingHost
		}
		if r.TTL == 0 {
			return errMissingTTL
		}

		// check values
		if !reSubDomain.MatchString(r.Name) {
			return errInvalidName
		}
		if !reTopDomain.MatchString(r.Value) {
			return errInvalidHost
		}
		if r.TTL < 1 {
			return errInvalidTTL
		}

		return nil
	case RecordTypeMX:
		// check for required attributes
		if r.Name == "" {
			return errMissingName
		}
		if r.Value == "" {
			return errMissingHost
		}
		if r.TTL == 0 {
			return errMissingTTL
		}
		if r.Priority.Valid == false {
			return errMissingPriority
		}

		// check values
		if !reSubDomain.MatchString(r.Name) {
			return errInvalidName
		}
		if !reMXDomain.MatchString(r.Value) {
			return errInvalidHost
		}
		if reIPv4Address.MatchString(r.Value) {
			return errInvalidHost
		}
		if r.TTL < 1 {
			return errInvalidTTL
		}
		if r.Priority.Int32 < 1 {
			return errInvalidPriority
		}

		return nil
	case RecordTypeNS:
		// check for required attributes
		if r.Name == "" {
			return errMissingName
		}
		if r.Value == "" {
			return errMissingHost
		}
		if r.TTL == 0 {
			return errMissingTTL
		}

		// check values
		if !reNSDomain.MatchString(r.Name) {
			return errInvalidName
		}
		if !reTopDomain.MatchString(r.Value) {
			return errInvalidHost
		}
		if r.TTL < 1 {
			return errInvalidTTL
		}

		return nil
	case RecordTypeTXT:
		// check for required attributes
		if r.Name == "" {
			return errMissingName
		}
		if r.Value == "" {
			return errMissingText
		}
		if r.TTL == 0 {
			return errMissingTTL
		}

		// check values
		if !reSubDomain.MatchString(r.Name) {
			return errInvalidName
		}
		if len(r.Value) > 255 {
			return errLengthExceededText
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
