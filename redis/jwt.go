package redis

import (
	redisCon "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"time"
)

// DeleteAccessToken deletes an access token from redis.
func (c *Client) DeleteAccessToken(accessTokenId uuid.UUID) (int, error) {
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
	reply, err = conn.Do("DEL", keyJwtAccess(accessTokenId.String()))
	if err != nil {
		return 0, err
	}

	return redisCon.Int(reply, nil)
}

// DeleteRefreshToken deletes a refresh token from redis.
func (c *Client) DeleteRefreshToken(refreshTokenId string) (int, error) {
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
	reply, err = conn.Do("DEL", keyJwtRefresh(refreshTokenId))
	if err != nil {
		return 0, err
	}

	return redisCon.Int(reply, nil)
}

// GetAccessToken retrieves an access token from redis.
func (c *Client) GetAccessToken(accessTokenId uuid.UUID) (uuid.UUID, error) {
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

	reply, err = conn.Do("GET", keyJwtAccess(accessTokenId.String()))
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
func (c *Client) SetAccessToken(accessTokenId, userId uuid.UUID, expire time.Duration) error {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		logger.Errorf("error connecting to redis")
		return errCantConnect
	}
	defer conn.Close()

	// add key
	var err error
	_, err = conn.Do("SET", keyJwtAccess(accessTokenId.String()), userId.String(), "EX", int(expire.Seconds()))
	if err != nil {
		return err
	}

	return nil
}

// SetRefreshToken adds a refresh token to redis.
func (c *Client) SetRefreshToken(refreshTokenId string, userId uuid.UUID, expire time.Duration) error {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		logger.Errorf("error connecting to redis")
		return errCantConnect
	}
	defer conn.Close()

	// add key
	var err error
	_, err = conn.Do("SET", keyJwtRefresh(refreshTokenId), userId.String(), "EX", int(expire.Seconds()))
	if err != nil {
		return err
	}

	return nil
}
