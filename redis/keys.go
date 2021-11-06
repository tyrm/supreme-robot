package redis

var (
	KeyBase = "sr:"

	KeySession = KeyBase + "session:"
	KeyDomains = KeyBase + "domains"
	KeyZones   = KeyBase + "dns:"

	KeyJwt          = KeyBase + "jwt:"
	KeyJwtAccesses  = KeyJwt + "a:"
	KeyJwtRefreshes = KeyJwt + "r:"
)

func KeyZone(d string) string       { return KeyZones + d }
func KeyJwtAccess(d string) string  { return KeyJwtAccesses + d }
func KeyJwtRefresh(d string) string { return KeyJwtRefreshes + d }
