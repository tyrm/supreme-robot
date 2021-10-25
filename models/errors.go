package models

import "errors"

var (
	ErrAlreadyCreated   = errors.New("model already created")
	ErrNotCreated       = errors.New("not created")
	ErrUnknownAttribute = errors.New("unknown attribute")
)
