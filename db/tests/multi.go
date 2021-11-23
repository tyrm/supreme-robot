package tests

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"testing"
)

// DoCreateDomainWRecords tests the CreateDomainWRecords function
func DoCreateDomainWRecords(t *testing.T, client db.DB) {
	newDomain := models.Domain{
		Domain:  "docreatedomainwrecords.",
		OwnerID: userAdmin.ID,
	}
	newRecordSOA := models.Record{
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
	err := newRecordSOA.Validate()
	if err != nil {
		t.Errorf("invalid soa record: %s", err.Error())
	}
	newRecordNS := models.Record{
		Name:  "@",
		Type:  models.RecordTypeNS,
		Value: "ns1.example.com",
		TTL:   300,
	}

	err = client.CreateDomainWRecords(&newDomain, &newRecordSOA, &newRecordNS)
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
	if newRecordSOA.ID == uuid.Nil {
		t.Errorf("record id not set, got: %s, want not: %s", newRecordSOA.ID.String(), uuid.Nil.String())
	}

	if newRecordSOA.DomainID == newDomain.OwnerID {
		t.Errorf("record id not set, got: %s, want not: %s", newRecordSOA.DomainID.String(), newDomain.OwnerID.String())
	}

	if newRecordSOA.CreatedAt.Unix() <= 0 {
		t.Errorf("record CreatedAt not set, got: %d, want not: >0", newRecordSOA.CreatedAt.Unix())
	}

	if newRecordSOA.UpdatedAt.Unix() <= 0 {
		t.Errorf("record UpdatedAt not set, got: %d, want not: >0", newRecordSOA.UpdatedAt.Unix())
	}
}

// DoCreateGroupsForUser tests the CreateGroupsForUser function for a known user
func DoCreateGroupsForUser(t *testing.T, client db.DB) {
	newUser := models.User{
		Username: "docreategroupforuser",
	}
	err := newUser.SetPassword("newpassword")
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

// DoCreateGroupsForUserUnknownUser tests the CreateGroupsForUser for an unknown user
func DoCreateGroupsForUserUnknownUser(t *testing.T, client db.DB) {
	id := uuid.MustParse("241a4327-a0ca-41df-855d-e0ecf552802c")
	err := client.CreateGroupsForUser(id, models.GroupDNSAdmin, models.GroupUserAdmin)
	if err.Error() != db.ErrUnknownUser.Error() {
		t.Fatalf("unexpected error, got: %s, want: %s", err.Error(), db.ErrUnknownUser.Error())
	}
}
