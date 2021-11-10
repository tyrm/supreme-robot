package db

import "errors"

var (
	ErrUnknownAttribute = errors.New("unknown attribute")
	ErrUnknownType      = errors.New("unknown type")
)
