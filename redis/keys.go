package redis

import "fmt"

var (
	KeyBase = "sr:"

	KeySession = KeyBase + "session:"
	KeyDomains = KeyBase + "domains"
	KeyZones = KeyBase + "dns:"
)

func KeyZone(d string) string { return fmt.Sprintf("%s%s", KeyZones, d)}
