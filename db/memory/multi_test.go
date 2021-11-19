package memory

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"testing"
)

func TestClient_CreateDomainWRecords(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	newDomain := models.Domain{
		Domain: "test.",
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

	err = client.CreateDomainWRecords(&newDomain, &newRecord)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	// Validate Domain
	if newDomain.ID == uuid.Nil {
		t.Errorf("domain id not set, got: %s, want not: %s", newDomain.ID.String(), uuid.Nil.String())
	}

	if newDomain.CreatedAt.Unix() <= 0 {
		t.Errorf("domain CreatedAt not set, got: %d, want not: >0", newDomain.CreatedAt.Unix())
	}

	if newDomain.UpdatedAt.Unix() <= 0 {
		t.Errorf("domain UpdatedAt not set, got: %d, want not: >0", newDomain.UpdatedAt.Unix())
	}

	// Validate Record
	if newRecord.ID == uuid.Nil {
		t.Errorf("record id not set, got: %s, want not: %s", newRecord.ID.String(), uuid.Nil.String())
	}

	if newRecord.DomainID == newDomain.OwnerID {
		t.Errorf("record id not set, got: %s, want not: %s", newRecord.DomainID.String(), newDomain.OwnerID.String())
	}

	if newRecord.CreatedAt.Unix() <= 0 {
		t.Errorf("record CreatedAt not set, got: %d, want not: >0", newRecord.CreatedAt.Unix())
	}

	if newRecord.UpdatedAt.Unix() <= 0 {
		t.Errorf("record UpdatedAt not set, got: %d, want not: >0", newRecord.UpdatedAt.Unix())
	}
}

func TestClient_CreateGroupsForUser(t *testing.T) {
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

	err = client.CreateGroupsForUser(newUser.ID, models.GroupDNSAdmin, models.GroupUserAdmin)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	user2, err := client.ReadUser(newUser.ID)

	searchGroups := map[uuid.UUID]bool{
		models.GroupDNSAdmin:  false,
		models.GroupUserAdmin: false,
	}
	for _, g := range user2.Groups {
		for k := range searchGroups {
			if g == k {
				searchGroups[g] = true
			}
		}
	}
	for k, v := range searchGroups {
		if !v {
			t.Errorf("didn't find group: %s", k)
		}
	}
}

func TestClient_CreateGroupsForUser_UnknownUser(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	id := uuid.MustParse("241a4327-a0ca-41df-855d-e0ecf552802c")
	err = client.CreateGroupsForUser(id, models.GroupDNSAdmin, models.GroupUserAdmin)
	if err.Error() != db.ErrUnknownUser.Error() {
		t.Fatalf("unexpected error, got: %s, want: %s", err.Error(), db.ErrUnknownUser.Error())
	}

}
