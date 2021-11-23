package tests

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"testing"
)

// DoReadRecordsForDomainOrderByName tests the CreateDomainWRecords function ordering by name
func DoReadRecordsForDomainOrderByName(t *testing.T, client db.DB) {
	newDomain := models.Domain{
		Domain:  "doreadrecorsfordomainorderbyname.",
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
	newRecordNS := models.Record{
		Name:  "@",
		Type:  models.RecordTypeNS,
		Value: "ns1.example.com.",
		TTL:   300,
	}

	err := client.CreateDomainWRecords(&newDomain, &newRecordSOA, &newRecordNS)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	receivedRecords, err := client.ReadRecordsForDomain(newDomain.ID, "name", true)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}
	receivedRecordsCount := len(*receivedRecords)
	if receivedRecordsCount != 2 {
		t.Errorf("invalid number of records returned, got: %d, want: 2", receivedRecordsCount)
	}

	// check records
	searchRecords := []models.Record{
		newRecordSOA,
		newRecordNS,
	}

	foundRecords := make(map[uuid.UUID]bool)
	for _, d := range searchRecords {
		foundRecords[d.ID] = false
	}

	for _, ar := range *receivedRecords {
		for _, sr := range searchRecords {
			if ar.ID == sr.ID {
				nameMatch := ar.Name == sr.Name
				domainIDMatch := ar.DomainID == sr.DomainID
				typeMatch := ar.Type == sr.Type
				valueMatch := ar.Value == sr.Value
				createdAtMatch := ar.CreatedAt == sr.CreatedAt
				updatedAtMatch := ar.UpdatedAt == sr.UpdatedAt

				if nameMatch && domainIDMatch && typeMatch && valueMatch && createdAtMatch && updatedAtMatch {
					foundRecords[ar.ID] = true
				}
			}
		}
	}

	for k, v := range foundRecords {
		if !v {
			t.Errorf("didn't find expected record: %s", k)
		}
	}
}

// DoReadRecordsForDomainOrderByUnknown tests the ReadRecordsForDomain for an unknown sort type
func DoReadRecordsForDomainOrderByUnknown(t *testing.T, client db.DB) {
	newDomain := models.Domain{
		Domain:  "doreadrecorsfordomainorderbyunknown.",
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
	newRecordNS := models.Record{
		Name:  "@",
		Type:  models.RecordTypeNS,
		Value: "ns1.example.com",
		TTL:   300,
	}

	err := client.CreateDomainWRecords(&newDomain, &newRecordSOA, &newRecordNS)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	receivedRecords, err := client.ReadRecordsForDomain(newDomain.ID, "unknown", true)
	if err != db.ErrUnknownAttribute {
		t.Errorf("unexpected error, got: %s, want: %s", err.Error(), db.ErrUnknownAttribute)
	}
	if receivedRecords != nil {
		t.Errorf("unexpected records, got: %#v, want: nil", receivedRecords)
	}
}
