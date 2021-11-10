package models

import "errors"

var (
	errInvalidIP          = errors.New("invalid ip")
	errInvalidHost        = errors.New("invalid host")
	errInvalidName        = errors.New("invalid name")
	errInvalidPort        = errors.New("invalid port")
	errInvalidPriority    = errors.New("invalid priority")
	errInvalidTTL         = errors.New("invalid ttl")
	errInvalidWeight      = errors.New("invalid weight")
	errLengthExceededText = errors.New("text length exceeded")
	errMissingIP          = errors.New("ip must be defined")
	errMissingHost        = errors.New("host must be defined")
	errMissingName        = errors.New("name must be defined")
	errMissingPort        = errors.New("port must be defined")
	errMissingPriority    = errors.New("priority must be defined")
	errMissingText        = errors.New("text must be defined")
	errMissingTTL         = errors.New("ttl must be defined")
	errMissingType        = errors.New("type must be defined")
	errMissingWeight      = errors.New("weight must be defined")
	errUnknownType        = errors.New("unknown type")
)
