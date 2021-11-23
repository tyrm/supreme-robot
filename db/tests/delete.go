package tests

import (
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"testing"
)

// DoDeleteDomain tests the Delete function for a models.Domain type
func DoDeleteDomain(t *testing.T, client db.DB) {
	newDomain := models.Domain{
		Domain: "dodeletedomain.",
	}
	err := client.Create(&newDomain)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	err = client.Delete(&newDomain)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	receivedDomain, err := client.ReadDomain(newDomain.ID)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}
	if receivedDomain != nil {
		t.Errorf("unexpected domain, got: %v, want: nil.", receivedDomain)
	}

	receivedDomainZ, err := client.ReadDomainZ(newDomain.ID)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}
	if receivedDomainZ == nil {
		t.Errorf("unexpected nil, got: nil, want: models.Domain")
	}
}

// DoDeleteUnknownType tests Delete for an unknown type
func DoDeleteUnknownType(t *testing.T, client db.DB) {
	newUnknown := unknownType{}
	err := client.Delete(&newUnknown)
	if err.Error() != db.ErrUnknownType.Error() {
		t.Fatalf("unexpected error, got: %s, want: %s", err.Error(), db.ErrUnknownType.Error())
	}
}
