package memory

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"testing"
)

type unknownType struct {
}

func TestClient_Create_Domain(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	newDomain := models.Domain{
		Domain: "test.",
	}
	err = client.Create(&newDomain)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	if newDomain.ID == uuid.Nil {
		t.Errorf("domain id not set, got: %s, want not: %s", newDomain.ID.String(), uuid.Nil.String())
	}

	if newDomain.CreatedAt.Unix() <= 0 {
		t.Errorf("domain CreatedAt not set, got: %d, want not: >0", newDomain.CreatedAt.Unix())
	}

	if newDomain.UpdatedAt.Unix() <= 0 {
		t.Errorf("domain UpdatedAt not set, got: %d, want not: >0", newDomain.UpdatedAt.Unix())
	}
}

func TestClient_Create_Record(t *testing.T) {
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

	newRecord := models.Record{
		Name:  "@",
		Type:  models.RecordTypeSOA,
		Value: "ns.example.com.",
		TTL:   300,
		MBox: sql.NullString{
			String: "hostmaster.example.com.",
			Valid:  true,
		},
		Refresh: sql.NullInt32{
			Int32: 22,
			Valid: true,
		},
		Retry: sql.NullInt32{
			Int32: 44,
			Valid: true,
		},
		Expire: sql.NullInt32{
			Int32: 33,
			Valid: true,
		},
	}

	err = newRecord.Validate()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	err = client.Create(&newRecord)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	if newRecord.ID == uuid.Nil {
		t.Errorf("record id not set, got: %s, want not: %s", newRecord.ID.String(), uuid.Nil.String())
	}

	if newRecord.CreatedAt.Unix() <= 0 {
		t.Errorf("record CreatedAt not set, got: %d, want not: >0", newRecord.CreatedAt.Unix())
	}

	if newRecord.UpdatedAt.Unix() <= 0 {
		t.Errorf("record UpdatedAt not set, got: %d, want not: >0", newRecord.UpdatedAt.Unix())
	}
}

func TestClient_Create_User(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	newUser := models.User{
		Username: "newuser",
	}
	err = newUser.SetPassword("newpassword")
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}
	err = client.Create(&newUser)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	if newUser.ID == uuid.Nil {
		t.Errorf("domain id not set, got: %s, want not: %s", newUser.ID.String(), uuid.Nil.String())
	}

	if newUser.CreatedAt.Unix() <= 0 {
		t.Errorf("domain CreatedAt not set, got: %d, want not: >0", newUser.CreatedAt.Unix())
	}

	if newUser.UpdatedAt.Unix() <= 0 {
		t.Errorf("domain UpdatedAt not set, got: %d, want not: >0", newUser.UpdatedAt.Unix())
	}
}

func TestClient_Create_UnknownType(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	newUnknown := unknownType{}
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}
	err = client.Create(&newUnknown)
	if err.Error() != db.ErrUnknownType.Error() {
		t.Fatalf("unexpected error, got: %s, want: %s", err.Error(), db.ErrUnknownType.Error())
	}
}
