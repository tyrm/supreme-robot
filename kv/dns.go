package kv

// DNS is for updating data for CoreDNS
type DNS interface {
	AddDomain(d string) error
	GetDomains() (*[]string, error)
	RemoveDomain(d string) error
}
