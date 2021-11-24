package memory

import (
	"github.com/patrickmn/go-cache"
	"github.com/tyrm/supreme-robot/kv"
	"github.com/tyrm/supreme-robot/util"
)

// AddDomain will add a domain name to the list of domains.
func (c *Client) AddDomain(d string) error {
	c.Lock()
	defer c.Unlock()

	if domains, domainsFound := c.KV.Get(kv.KeyDomains()); domainsFound {
		domainList := domains.([]string)
		domainList = append(domainList, d)
		c.KV.Set(kv.KeyDomains(), domainList, cache.NoExpiration)
		return nil
	}

	newList := []string{d}

	c.KV.Set(kv.KeyDomains(), newList, cache.NoExpiration)
	return nil
}

// RemoveDomain will remove a domain name from the list of domains.
func (c *Client) RemoveDomain(d string) error {
	c.Lock()
	defer c.Unlock()

	if domains, domainsFound := c.KV.Get(kv.KeyDomains()); domainsFound {
		domainList := domains.([]string)
		domainList = util.FastPopString(domainList, d)
		c.KV.Set(kv.KeyDomains(), domainList, cache.NoExpiration)
		return nil
	}

	return nil
}

// GetDomains returns all domains.
func (c *Client) GetDomains() (*[]string, error) {
	c.RLock()
	defer c.RUnlock()

	if domains, domainsFound := c.KV.Get(kv.KeyDomains()); domainsFound {
		domainList := domains.([]string)
		return &domainList, nil
	}

	return nil, nil
}
