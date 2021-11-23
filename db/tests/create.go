package tests

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"testing"
)

type unknownType struct {
}

// DoCreateDomain tests the Create function for a models.Domain type
func DoCreateDomain(t *testing.T, client db.DB) {
	newDomain := models.Domain{
		Domain:  "docreatedomain.",
		OwnerID: userAdmin.ID,
	}
	err := client.Create(&newDomain)
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

// DoCreateRecord tests the Create function for a models.Record type
func DoCreateRecord(t *testing.T, client db.DB) {
	newDomain := models.Domain{
		Domain:  "docreaterecord.",
		OwnerID: userAdmin.ID,
	}
	err := client.Create(&newDomain)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	newRecord := models.Record{
		DomainID: newDomain.ID,
		Name:     "@",
		Type:     models.RecordTypeSOA,
		Value:    "ns.example.com.",
		TTL:      300,
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

// DoCreateUser tests the Create function for a models.User type
func DoCreateUser(t *testing.T, client db.DB) {
	newUser := models.User{
		Username: "newuser",
	}
	err := newUser.SetPassword("newpassword")
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

// DoCreateUnknownType tests Create for an unknown type
func DoCreateUnknownType(t *testing.T, client db.DB) {
	newUnknown := unknownType{}
	err := client.Create(&newUnknown)
	if err.Error() != db.ErrUnknownType.Error() {
		t.Fatalf("unexpected error, got: %s, want: %s", err.Error(), db.ErrUnknownType.Error())
	}
}
