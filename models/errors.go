package models

import "errors"

var (
	ErrAlreadyCreated = errors.New("model already created")
	ErrUnknownAttribute  = errors.New("unknown attribute")
)