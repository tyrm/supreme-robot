package redis

import (
	redisCon "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/kv"
	"time"
)

// DeleteAccessToken deletes an access token from redis.
func (c *Client) DeleteAccessToken(accessTokenID uuid.UUID) (int, error) {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		logger.Errorf("error connecting to redis")
		return 0, errCantConnect
	}
	defer conn.Close()

	// add key
	var (
		err   error
		reply interface{}
	)
	reply, err = conn.Do("DEL", kv.KeyJwtAccess(accessTokenID.String()))
	if err != nil {
		return 0, err
	}

	return redisCon.Int(reply, nil)
}

// DeleteRefreshToken deletes a refresh token from redis.
func (c *Client) DeleteRefreshToken(refreshTokenID string) (int, error) {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		logger.Errorf("error connecting to redis")
		return 0, errCantConnect
	}
	defer conn.Close()

	// add key
	var (
		err   error
		reply interface{}
	)
	reply, err = conn.Do("DEL", kv.KeyJwtRefresh(refreshTokenID))
	if err != nil {
		return 0, err
	}

	return redisCon.Int(reply, nil)
}

// GetAccessToken retrieves an access token from redis.
func (c *Client) GetAccessToken(accessTokenID uuid.UUID) (uuid.UUID, error) {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		logger.Errorf("error connecting to redis")
		return uuid.Nil, errCantConnect
	}
	defer conn.Close()

	var (
		reply interface{}
		err   error
		val   string
	)

	reply, err = conn.Do("GET", kv.KeyJwtAccess(accessTokenID.String()))
	if err != nil {
		return uuid.Nil, err
	}

	val, err = redisCon.String(reply, nil)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(val)
}

// SetAccessToken adds an access token to redis.
func (c *Client) SetAccessToken(accessTokenID, userID uuid.UUID, expire time.Duration) error {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		logger.Errorf("error connecting to redis")
		return errCantConnect
	}
	defer conn.Close()

	// add key
	var err error
	_, err = conn.Do("SET", kv.KeyJwtAccess(accessTokenID.String()), userID.String(), "EX", int(expire.Seconds()))
	if err != nil {
		return err
	}

	return nil
}

// SetRefreshToken adds a refresh token to redis.
func (c *Client) SetRefreshToken(refreshTokenID string, userID uuid.UUID, expire time.Duration) error {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		logger.Errorf("error connecting to redis")
		return errCantConnect
	}
	defer conn.Close()

	// add key
	var err error
	_, err = conn.Do("SET", kv.KeyJwtRefresh(refreshTokenID), userID.String(), "EX", int(expire.Seconds()))
	if err != nil {
		return err
	}

	return nil
}
