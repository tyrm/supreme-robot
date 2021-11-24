package tests

import (
	"database/sql"
	"fmt"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"testing"
	"time"
)

// DoReadRecordsForDomain tests the CreateDomainWRecords function ordering by name
func DoReadRecordsForDomain(t *testing.T, client db.DB) {
	// prep data for test
	newDomain := models.Domain{
		Domain:  "doreadrecorsfordomainorderbyname.",
		OwnerID: userAdmin.ID,
	}
	err := client.Create(&newDomain)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	newRecordSOA := models.Record{
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
	newRecordNS := models.Record{
		DomainID: newDomain.ID,
		Name:     "@",
		Type:     models.RecordTypeNS,
		Value:    "ns1.example.com.",
		TTL:      300,
	}
	newRecordAAlice := models.Record{
		DomainID: newDomain.ID,
		Name:     "alice",
		Type:     models.RecordTypeA,
		Value:    "1.2.3.4",
		TTL:      300,
	}
	newRecordABill := models.Record{
		DomainID: newDomain.ID,
		Name:     "bill",
		Type:     models.RecordTypeA,
		Value:    "4.5.6.7",
		TTL:      300,
	}
	newRecordACharlie := models.Record{
		DomainID: newDomain.ID,
		Name:     "charlie",
		Type:     models.RecordTypeA,
		Value:    "8.9.1.2",
		TTL:      300,
	}

	insertOrder := []*models.Record{
		&newRecordABill,
		&newRecordNS,
		&newRecordACharlie,
		&newRecordSOA,
		&newRecordAAlice,
	}
	for _, r := range insertOrder {
		err = client.Create(r)
		if err != nil {
			t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
			return
		}
		time.Sleep(1 * time.Second)
	}

	nameAscOrder := make([]*models.Record, 5)
	nameDescOrder := make([]*models.Record, 5)

	if newRecordSOA.ID.String() < newRecordNS.ID.String() {
		nameAscOrder = []*models.Record{&newRecordSOA, &newRecordNS, &newRecordAAlice, &newRecordABill, &newRecordACharlie}
		nameDescOrder = []*models.Record{&newRecordACharlie, &newRecordABill, &newRecordAAlice, &newRecordNS, &newRecordSOA}
	} else {
		nameAscOrder = []*models.Record{&newRecordNS, &newRecordSOA, &newRecordAAlice, &newRecordABill, &newRecordACharlie}
		nameDescOrder = []*models.Record{&newRecordACharlie, &newRecordABill, &newRecordAAlice, &newRecordSOA, &newRecordNS}
	}

	// do the tests
	tables := []struct {
		orderBy       string
		asc           bool
		expectedOrder []*models.Record
	}{
		{"name", true, nameAscOrder},
		{"name", false, nameDescOrder},
		{"created_at", true, []*models.Record{&newRecordABill, &newRecordNS, &newRecordACharlie, &newRecordSOA, &newRecordAAlice}},
		{"created_at", false, []*models.Record{&newRecordAAlice, &newRecordSOA, &newRecordACharlie, &newRecordNS, &newRecordABill}},
	}

	for _, table := range tables {
		table := table
		name := fmt.Sprintf("[%s,%v] Test", table.orderBy, table.asc)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			receivedRecords, err := client.ReadRecordsForDomain(newDomain.ID, table.orderBy, table.asc)
			if err != nil {
				t.Errorf("[%s,%v] unexpected error, got: %s, want: nil", table.orderBy, table.asc, err.Error())
				return
			}
			receivedRecordsCount := len(*receivedRecords)
			if receivedRecordsCount != 5 {
				t.Errorf("[%s,%v] invalid number of records returned, got: %d, want: 5", table.orderBy, table.asc, receivedRecordsCount)
				return
			}

			// check records
			for i, ar := range *receivedRecords {
				if ar.ID != table.expectedOrder[i].ID {
					t.Errorf("[%s,%v][%d] unexpected ID, got: %s, want: %s", table.orderBy, table.asc, i, ar.ID, table.expectedOrder[i].ID)
				}
				if ar.Name != table.expectedOrder[i].Name {
					t.Errorf("[%s,%v][%d] unexpected Name, got: %s, want: %s", table.orderBy, table.asc, i, ar.Name, table.expectedOrder[i].Name)
				}
				if ar.Type != table.expectedOrder[i].Type {
					t.Errorf("[%s,%v][%d] unexpected Type, got: %s, want: %s", table.orderBy, table.asc, i, ar.Type, table.expectedOrder[i].Type)
				}
				if ar.Value != table.expectedOrder[i].Value {
					t.Errorf("[%s,%v][%d] unexpected Value, got: %s, want: %s", table.orderBy, table.asc, i, ar.Value, table.expectedOrder[i].Value)
				}
			}

		})
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

// DoReadRecordsForDomainByName tests the ReadRecordsForDomainByName function
func DoReadRecordsForDomainByName(t *testing.T, client db.DB) {
	// prep data for test
	newDomain := models.Domain{
		Domain:  "doreadrecordsfordomainbyname.",
		OwnerID: userAdmin.ID,
	}
	err := client.Create(&newDomain)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	newRecordSOA := models.Record{
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
	newRecordNS := models.Record{
		DomainID: newDomain.ID,
		Name:     "@",
		Type:     models.RecordTypeNS,
		Value:    "ns1.example.com.",
		TTL:      300,
	}
	newRecordAAlice := models.Record{
		DomainID: newDomain.ID,
		Name:     "alice",
		Type:     models.RecordTypeA,
		Value:    "1.2.3.4",
		TTL:      300,
	}
	newRecordABill := models.Record{
		DomainID: newDomain.ID,
		Name:     "bill",
		Type:     models.RecordTypeA,
		Value:    "4.5.6.7",
		TTL:      300,
	}
	newRecordACharlie := models.Record{
		DomainID: newDomain.ID,
		Name:     "charlie",
		Type:     models.RecordTypeA,
		Value:    "8.9.1.2",
		TTL:      300,
	}

	insertOrder := []*models.Record{
		&newRecordABill,
		&newRecordNS,
		&newRecordACharlie,
		&newRecordSOA,
		&newRecordAAlice,
	}
	for _, r := range insertOrder {
		err = client.Create(r)
		if err != nil {
			t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
			return
		}
		time.Sleep(1 * time.Second)
	}

	// do the tests
	tables := []struct {
		name          string
		expectedOrder []*models.Record
	}{
		{"@", []*models.Record{&newRecordNS, &newRecordSOA}},
		{"alice", []*models.Record{&newRecordAAlice}},
		{"bill", []*models.Record{&newRecordABill}},
		{"charlie", []*models.Record{&newRecordACharlie}},
	}

	for _, table := range tables {
		table := table
		name := fmt.Sprintf("Testing %s", table.name)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			receivedRecords, err := client.ReadRecordsForDomainByName(newDomain.ID, table.name)
			if err != nil {
				t.Errorf("[%s] unexpected error, got: %s, want: nil", table.name, err.Error())
				return
			}
			receivedRecordsCount := len(*receivedRecords)
			expectedRecordsCount := len(table.expectedOrder)
			if receivedRecordsCount != expectedRecordsCount {
				t.Errorf("[%s] invalid number of records returned, got: %d, want: %d", table.name, receivedRecordsCount, expectedRecordsCount)
				return
			}

			// check records
			for i, ar := range *receivedRecords {
				if ar.ID != table.expectedOrder[i].ID {
					t.Errorf("[%s][%d] unexpected ID, got: %s, want: %s", table.name, i, ar.ID, table.expectedOrder[i].ID)
				}
				if ar.Name != table.expectedOrder[i].Name {
					t.Errorf("[%s][%d] unexpected Name, got: %s, want: %s", table.name, i, ar.Name, table.expectedOrder[i].Name)
				}
				if ar.Type != table.expectedOrder[i].Type {
					t.Errorf("[%s][%d] unexpected Type, got: %s, want: %s", table.name, i, ar.Type, table.expectedOrder[i].Type)
				}
				if ar.Value != table.expectedOrder[i].Value {
					t.Errorf("[%s][%d] unexpected Value, got: %s, want: %s", table.name, i, ar.Value, table.expectedOrder[i].Value)
				}
			}
		})
	}
}
