package redis

import (
	"fmt"
	redisCon "github.com/gomodule/redigo/redis"
)

// AddDomain will add a domain name to the list of domains.
func (c *Client) AddDomain(d string) error {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		fmt.Println("error connecting to redis")
		return errCantConnect
	}
	defer conn.Close()

	// add key
	var err error
	_, err = conn.Do("SADD", keyDomains, d)
	if err != nil {
		return err
	}

	return nil
}

// RemoveDomain will remove a domain name from the list of domains.
func (c *Client) RemoveDomain(d string) error {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		fmt.Println("error connecting to redis")
		return errCantConnect
	}
	defer conn.Close()

	// remove key
	var err error
	_, err = conn.Do("SREM", keyDomains, d)
	if err != nil {
		return err
	}

	return nil
}

// GetDomains returns all domains.
func (c *Client) GetDomains() (*[]string, error) {
	// get connection
	conn := c.db.Get()
	if conn == nil {
		fmt.Println("error connecting to redis")
		return nil, errCantConnect
	}
	defer conn.Close()

	var (
		reply interface{}
		err   error
		vals  []string
	)

	reply, err = conn.Do("SMEMBERS", keyDomains)
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
