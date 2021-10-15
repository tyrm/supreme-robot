package redis

import (
	"fmt"
	redisCon "github.com/gomodule/redigo/redis"
)

func (c *Client) AddDomain(d string) error {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		fmt.Println("error connecting to redis")
		return ErrCantConnect
	}
	defer conn.Close()

	// add key
	var err error
	_, err = conn.Do("SADD", KeyDomains, d)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveDomain(d string) error {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		fmt.Println("error connecting to redis")
		return ErrCantConnect
	}
	defer conn.Close()

	// remove key
	var err error
	_, err = conn.Do("SREM", KeyDomains, d)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetDomains() (*[]string, error) {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		fmt.Println("error connecting to redis")
		return nil, ErrCantConnect
	}
	defer conn.Close()

	var (
		reply interface{}
		err error
		vals []string
	)

	reply, err = conn.Do("SMEMBERS", KeyDomains)
	if err != nil {
		return nil, err
	}

	vals, err = redisCon.Strings(reply, nil)
	if err != nil {
		return nil, err
	}

	logger.Tracef("domains: %v", vals)
	return &vals, nil
}