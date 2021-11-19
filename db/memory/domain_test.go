package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/models"
	"testing"
)

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
