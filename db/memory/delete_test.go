package memory

import (
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"testing"
)

func TestClient_Delete_Domain(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	newDomain := models.Domain{
		Domain: "test.",
	}
	err = client.Create(&newDomain)
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

func TestClient_Delete_UnknownType(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	newUnknown := unknownType{}
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}
	err = client.Delete(&newUnknown)
	if err.Error() != db.ErrUnknownType.Error() {
		t.Fatalf("unexpected error, got: %s, want: %s", err.Error(), db.ErrUnknownType.Error())
	}
}
