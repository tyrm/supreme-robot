package models

import "errors"

var (
	errInvalidIP   = errors.New("invalid ip")
	errInvalidHost = errors.New("invalid host")
	errInvalidName = errors.New("invalid name")
	errInvalidTTL  = errors.New("invalid ttl")
	errMissingIP   = errors.New("ip must be defined")
	errMissingHost = errors.New("host must be defined")
	errMissingName = errors.New("name must be defined")
	errMissingTTL  = errors.New("ttl must be defined")
	errMissingType = errors.New("type must be defined")
	errUnknownType = errors.New("unknown type")
)
