package kv

import (
	"github.com/google/uuid"
	"time"
)

type Webapp interface {
	DeleteAccessToken(uuid.UUID) (int, error)
	DeleteRefreshToken(string) (int, error)
	GetAccessToken(uuid.UUID) (uuid.UUID, error)
	SetAccessToken(uuid.UUID, uuid.UUID, time.Duration) error
	SetRefreshToken(string, uuid.UUID, time.Duration) error
}
