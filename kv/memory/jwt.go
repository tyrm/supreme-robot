package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/kv"
	"time"
)

// DeleteAccessToken deletes an access token from redis.
func (c *Client) DeleteAccessToken(accessTokenID uuid.UUID) (int, error) {
	c.KV.Delete(kv.KeyJwtAccess(accessTokenID.String()))
	return 1, nil
}

// DeleteRefreshToken deletes a refresh token from redis.
func (c *Client) DeleteRefreshToken(refreshTokenID string) (int, error) {
	c.KV.Delete(kv.KeyJwtRefresh(refreshTokenID))
	return 1, nil
}

// GetAccessToken retrieves an access token from redis.
func (c *Client) GetAccessToken(accessTokenID uuid.UUID) (uuid.UUID, error) {
	if userID, userIDFound := c.KV.Get(kv.KeyJwtAccess(accessTokenID.String())); userIDFound {
		return uuid.Parse(userID.(string))
	}
	return uuid.Nil, nil
}

// SetAccessToken adds an access token to redis.
func (c *Client) SetAccessToken(accessTokenID, userID uuid.UUID, expire time.Duration) error {
	c.KV.Set(kv.KeyJwtAccess(accessTokenID.String()), userID.String(), expire)
	return nil
}

// SetRefreshToken adds a refresh token to redis.
func (c *Client) SetRefreshToken(refreshTokenID string, userID uuid.UUID, expire time.Duration) error {
	c.KV.Set(kv.KeyJwtRefresh(refreshTokenID), userID.String(), expire)
	return nil
}
