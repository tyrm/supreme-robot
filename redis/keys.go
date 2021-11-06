package redis

var (
	keyBase = "sr:"

	keySession = keyBase + "session:"
	keyDomains = keyBase + "domains"
	keyZones   = keyBase + "dns:"

	keyJwt          = keyBase + "jwt:"
	keyJwtAccesses  = keyJwt + "a:"
	keyJwtRefreshes = keyJwt + "r:"
)

func keyZone(d string) string       { return keyZones + d }
func keyJwtAccess(d string) string  { return keyJwtAccesses + d }
func keyJwtRefresh(d string) string { return keyJwtRefreshes + d }
