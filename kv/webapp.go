package kv

import (
	"github.com/google/uuid"
	"time"
)

// Webapp is for dealing with kv values for use in the app
type Webapp interface {
	DeleteAccessToken(uuid.UUID) (int, error)
	DeleteRefreshToken(string) (int, error)
	GetAccessToken(uuid.UUID) (uuid.UUID, error)
	GetRefreshToken(string) (uuid.UUID, error)
	SetAccessToken(uuid.UUID, uuid.UUID, time.Duration) error
	SetRefreshToken(string, uuid.UUID, time.Duration) error
}
