package kv

var (
	keyBase = "sr:"

	keyDomains = keyBase + "domains"
	keyJwt     = keyBase + "jwt:"
	keyZones   = keyBase + "dns:"

	keyJwtAccesses  = keyJwt + "a:"
	keyJwtRefreshes = keyJwt + "r:"
)

// KeyDomains returns the kv key which holds a list of domains
func KeyDomains() string { return keyDomains }

// KeyJwtAccess returns the kv key which holds a JWT access token
func KeyJwtAccess(d string) string { return keyJwtAccesses + d }

// KeyJwtRefresh returns the kv key which holds a JWT refresh token
func KeyJwtRefresh(d string) string { return keyJwtRefreshes + d }

// KeyZone returns the kv key which holds the zone data for a domain
func KeyZone(d string) string { return keyZones + d }
