package kv

import "testing"

func TestKeyDomains(t *testing.T) {
	v := KeyDomains()
	if v != "sr:domains" {
		t.Errorf("enexpected value for KeyDomains, got: '%s', want: 'sr:domains'.", v)
	}
}

func TestKeyJwtAccess(t *testing.T) {
	v := KeyJwtAccess("test123")
	if v != "sr:jwt:a:test123" {
		t.Errorf("enexpected value for KeyDomains, got: '%s', want: 'sr:jwt:a:test123'.", v)
	}
}

func TestKeyJwtRefresh(t *testing.T) {
	v := KeyJwtRefresh("test123")
	if v != "sr:jwt:r:test123" {
		t.Errorf("enexpected value for KeyDomains, got: '%s', want: 'sr:jwt:r:test123'.", v)
	}
}

func TestKeyZone(t *testing.T) {
	v := KeyZone("example.com.")
	if v != "sr:dns:example.com." {
		t.Errorf("enexpected value for KeyDomains, got: '%s', want: 'sr:dns:example.com.'.", v)
	}
}
