package tests

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"testing"
)

// DoReadDomain tests the ReadDomain function for a known domain
func DoReadDomain(t *testing.T, client db.DB) {
	newDomain := models.Domain{
		Domain:  "docreatedomain.",
		OwnerID: userAdmin.ID,
	}
	err := client.Create(&newDomain)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	d, err := client.ReadDomain(newDomain.ID)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	// check records
	if d.ID != newDomain.ID {
		t.Errorf("unexpected id, got: %s, want: %s", d.ID.String(), newDomain.ID.String())
	}
	if d.Domain != newDomain.Domain {
		t.Errorf("unexpected domain, got: %s, want: %s", d.Domain, newDomain.Domain)
	}
}

// DoReadDomainUnknown tests the ReadDomain function for an unknown domain
func DoReadDomainUnknown(t *testing.T, client db.DB) {
	d, err := client.ReadDomain(uuid.MustParse("75062d21-3e64-4cb5-9109-0b87e785a8dc"))
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}
	if d != nil {
		t.Errorf("unexpected domain, got: %#v, want: nil", d)
	}
}

// DoReadDomainZ tests the ReadDomainZ function for a known domain
func DoReadDomainZ(t *testing.T, client db.DB) {
	newDomain := models.Domain{
		Domain:  "doreaddomainz.",
		OwnerID: userAdmin.ID,
	}
	err := client.Create(&newDomain)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	err = client.Delete(&newDomain)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	d, err := client.ReadDomainZ(newDomain.ID)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	// check records
	if d.ID != newDomain.ID {
		t.Errorf("unexpected id, got: %s, want: %s", d.ID.String(), newDomain.ID.String())
	}
	if d.Domain != newDomain.Domain {
		t.Errorf("unexpected domain, got: %s, want: %s", d.Domain, newDomain.Domain)
	}
}

// DoReadDomainZUnknown tests the ReadDomainZ function for an unknown domain
func DoReadDomainZUnknown(t *testing.T, client db.DB) {
	d, err := client.ReadDomainZ(uuid.MustParse("75062d21-3e64-4cb5-9109-0b87e785a8dc"))
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}
	if d != nil {
		t.Errorf("unexpected domain, got: %#v, want: nil", d)
	}
}

// DoReadDomainByDomain tests the ReadDomainByDomain function for a known domain
func DoReadDomainByDomain(t *testing.T, client db.DB) {
	newDomain := models.Domain{
		Domain:  "doreaddomainbydomain.",
		OwnerID: userAdmin.ID,
	}
	err := client.Create(&newDomain)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	d, err := client.ReadDomainByDomain(newDomain.Domain)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	// check records
	if d.ID != newDomain.ID {
		t.Errorf("unexpected id, got: %s, want: %s", d.ID.String(), newDomain.ID.String())
	}
	if d.Domain != newDomain.Domain {
		t.Errorf("unexpected domain, got: %s, want: %s", d.Domain, newDomain.Domain)
	}
}

// DoReadDomainByDomainUnknown tests the ReadDomainByDomain function for an unknown domain
func DoReadDomainByDomainUnknown(t *testing.T, client db.DB) {
	d, err := client.ReadDomainByDomain("meowmeow.")
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}
	if d != nil {
		t.Errorf("unexpected domain, got: %#v, want: nil", d)
	}
}

// DoReadDomainsForUserAdmin tests the ReadDomainsForUser function
func DoReadDomainsForUserAdmin(t *testing.T, client db.DB) {
	// create domains
	domains := []string{
		"doreaddomainsforuseradmin1.",
		"doreaddomainsforuseradmin2.",
		"doreaddomainsforuseradmin3.",
	}
	searchDomains := make([]models.Domain, 3)
	for i, d := range domains {
		newDomain := models.Domain{
			Domain:  d,
			OwnerID: userAdmin.ID,
		}
		err := client.Create(&newDomain)
		if err != nil {
			t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
			return
		}
		searchDomains[i] = newDomain
	}

	adminDomains, err := client.ReadDomainsForUser(userAdmin.ID)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}
	adminDomainsCount := len(*adminDomains)
	if adminDomainsCount != len(domains) {
		t.Errorf("invalid number of rows returned, got: %d, want: %d", adminDomainsCount, len(domains))
	}

	// check records
	foundDomains := make(map[uuid.UUID]bool)
	for _, d := range searchDomains {
		foundDomains[d.ID] = false
	}

	for _, ad := range *adminDomains {
		for _, sd := range searchDomains {
			if ad.ID == sd.ID {
				nameMatch := ad.Domain == sd.Domain
				ownerMatch := ad.OwnerID == sd.OwnerID
				createdAtMatch := ad.CreatedAt == sd.CreatedAt
				updatedAtMatch := ad.UpdatedAt == sd.UpdatedAt

				if nameMatch && ownerMatch && createdAtMatch && updatedAtMatch {
					foundDomains[ad.ID] = true
				}
			}
		}
	}

	for k, v := range foundDomains {
		if !v {
			t.Errorf("didn't find expected domain: %s", k)
		}
	}
}
