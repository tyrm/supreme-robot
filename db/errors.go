package db

import "errors"

var (
	// ErrUnknownAttribute returns when an invalid sorting attribute is selected
	ErrUnknownAttribute = errors.New("unknown attribute")
	// ErrUnknownType returns when an type is received by the client
	ErrUnknownType = errors.New("unknown type")
)
