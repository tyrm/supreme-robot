package postgres

import "errors"

var (
	errUnknownAttribute = errors.New("unknown attribute")
	errUnknownType      = errors.New("unknown type")
)