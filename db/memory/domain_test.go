package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/models"
	"testing"
)

func TestClient_ReadDomain(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	d, err := client.ReadDomain(domainExample.ID)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	// check records
	if d.ID != domainExample.ID {
		t.Errorf("unexpected id, got: %s, want: %s", d.ID.String(), domainExample.ID.String())
	}
	if d.Domain != domainExample.Domain {
		t.Errorf("unexpected domain, got: %s, want: %s", d.Domain, domainExample.Domain)
	}
}

func TestClient_ReadDomain_Unknown(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	d, err := client.ReadDomain(uuid.MustParse("75062d21-3e64-4cb5-9109-0b87e785a8dc"))
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}
	if d != nil {
		t.Errorf("unexpected domain, got: %#v, want: nil", d)
	}
}

func TestClient_ReadDomainZ(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	d, err := client.ReadDomainZ(domainExample.ID)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	// check records
	if d.ID != domainExample.ID {
		t.Errorf("unexpected id, got: %s, want: %s", d.ID.String(), domainExample.ID.String())
	}
	if d.Domain != domainExample.Domain {
		t.Errorf("unexpected domain, got: %s, want: %s", d.Domain, domainExample.Domain)
	}
}

func TestClient_ReadDomainZ_Unknown(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	d, err := client.ReadDomainZ(uuid.MustParse("75062d21-3e64-4cb5-9109-0b87e785a8dc"))
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}
	if d != nil {
		t.Errorf("unexpected domain, got: %#v, want: nil", d)
	}
}

func TestClient_ReadDomainByDomain(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	d, err := client.ReadDomainByDomain(domainTest3.Domain)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	// check records
	if d.ID != domainTest3.ID {
		t.Errorf("unexpected id, got: %s, want: %s", d.ID.String(), domainTest3.ID.String())
	}
	if d.Domain != domainTest3.Domain {
		t.Errorf("unexpected domain, got: %s, want: %s", d.Domain, domainTest3.Domain)
	}
}

func TestClient_ReadDomainByDomain_Unknown(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	d, err := client.ReadDomainByDomain("meowmeow.")
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}
	if d != nil {
		t.Errorf("unexpected domain, got: %#v, want: nil", d)
	}
}

func TestClient_ReadDomainsForUser_Admin(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	adminDomains, err := client.ReadDomainsForUser(uuid.MustParse("44892097-2c97-4c16-b4d1-e8522586df48"))
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}
	adminDomainsCount := len(*adminDomains)
	if adminDomainsCount != 3 {
		t.Errorf("invalid number of rows returned, got: %d, want: 3", adminDomainsCount)
	}

	// check records
	searchDomains := []models.Domain{
		domainTest1,
		domainTest2,
		domainTest3,
	}

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
