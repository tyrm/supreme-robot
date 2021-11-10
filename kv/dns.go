package kv

type DNS interface {
	AddDomain(d string) error
	RemoveDomain(d string) error
}
