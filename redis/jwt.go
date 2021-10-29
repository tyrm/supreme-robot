package redis

import (
	redisCon "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"time"
)

func (c *Client) DelAccessToken(accessTokenId uuid.UUID) (int, error) {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		logger.Errorf("error connecting to redis")
		return 0, ErrCantConnect
	}
	defer conn.Close()

	// add key
	var (
		err   error
		reply interface{}
	)
	reply, err = conn.Do("DEL", KeyJwtAccess(accessTokenId.String()))
	if err != nil {
		return 0, err
	}

	return redisCon.Int(reply, nil)
}

func (c *Client) DelRefreshToken(refreshTokenId string) (int, error) {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		logger.Errorf("error connecting to redis")
		return 0, ErrCantConnect
	}
	defer conn.Close()

	// add key
	var (
		err   error
		reply interface{}
	)
	reply, err = conn.Do("DEL", KeyJwtRefresh(refreshTokenId))
	if err != nil {
		return 0, err
	}

	return redisCon.Int(reply, nil)
}

func (c *Client) GetAccessToken(accessTokenId uuid.UUID) (uuid.UUID, error) {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		logger.Errorf("error connecting to redis")
		return uuid.Nil, ErrCantConnect
	}
	defer conn.Close()

	var (
		reply interface{}
		err   error
		val   string
	)

	reply, err = conn.Do("GET", KeyJwtAccess(accessTokenId.String()))
	if err != nil {
		return uuid.Nil, err
	}

	val, err = redisCon.String(reply, nil)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(val)
}

func (c *Client) SetAccessToken(accessTokenId, userId uuid.UUID, expire time.Duration) error {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		logger.Errorf("error connecting to redis")
		return ErrCantConnect
	}
	defer conn.Close()

	// add key
	var err error
	_, err = conn.Do("SET", KeyJwtAccess(accessTokenId.String()), userId.String(), "EX", int(expire.Seconds()))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SetRefreshToken(refreshTokenId string, userId uuid.UUID, expire time.Duration) error {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		logger.Errorf("error connecting to redis")
		return ErrCantConnect
	}
	defer conn.Close()

	// add key
	var err error
	_, err = conn.Do("SET", KeyJwtRefresh(refreshTokenId), userId.String(), "EX", int(expire.Seconds()))
	if err != nil {
		return err
	}

	return nil
}
