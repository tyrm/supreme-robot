package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/models"
	"testing"
)

func TestClient_ReadRecordsForDomain(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	receivedRecords, err := client.ReadRecordsForDomain(domainExample.ID, "name", true)
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
		recordExampleSOA,
		recordExampleNS1,
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
