package memory

import (
	"github.com/google/uuid"
	"time"
)

// DeleteAccessToken deletes an access token from redis.
func (c *Client) DeleteAccessToken(accessTokenID uuid.UUID) (int, error) {

	return 1, nil
}

// DeleteRefreshToken deletes a refresh token from redis.
func (c *Client) DeleteRefreshToken(refreshTokenID string) (int, error) {

	return 1, nil
}

// GetAccessToken retrieves an access token from redis.
func (c *Client) GetAccessToken(accessTokenID uuid.UUID) (uuid.UUID, error) {

	return uuid.Nil, nil
}

// SetAccessToken adds an access token to redis.
func (c *Client) SetAccessToken(accessTokenID, userID uuid.UUID, expire time.Duration) error {

	return nil
}

// SetRefreshToken adds a refresh token to redis.
func (c *Client) SetRefreshToken(refreshTokenID string, userID uuid.UUID, expire time.Duration) error {

	return nil
}
