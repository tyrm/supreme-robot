package models

import "errors"

var (
	errInvalidIP        = errors.New("invalid ip")
	errInvalidName      = errors.New("invalid name")
	errInvalidTTL       = errors.New("invalid ttl")
	errMissingIP        = errors.New("ip must be defined")
	errMissingName      = errors.New("name must be defined")
	errMissingTTL       = errors.New("ttl must be defined")
	errMissingType      = errors.New("type must be defined")
	errUnknownAttribute = errors.New("unknown attribute")
	errUnknownType      = errors.New("unknown type")
)
