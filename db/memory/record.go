package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"sort"
	"strings"
)

// ReadRecordsForDomain will read all records for a given domain
func (c *Client) ReadRecordsForDomain(domainID uuid.UUID, orderBy string, asc bool) (*[]models.Record, error) {
	records := make([]models.Record, 0)

	// Lock DB
	c.RLock()
	defer c.RUnlock()

	switch strings.ToLower(orderBy) {
	case "created_at":
	case "name":
		foundRecords := make(map[string]map[uuid.UUID]models.Record)

		// find records
		for _, record := range c.records {
			if record.DomainID == domainID {
				if foundRecords[record.Name] == nil {
					foundRecords[record.Name] = make(map[uuid.UUID]models.Record)
				}
				foundRecords[record.Name][record.ID] = record
			}
		}

		// get keys and sort
		var foundNames sort.StringSlice = make([]string, len(foundRecords))
		i := 0
		for k := range foundRecords {
			foundNames[i] = k
			i++
		}
		if asc {
			sort.Sort(foundNames)
		} else {
			sort.Sort(sort.Reverse(foundNames))
		}

		// make a list
		for _, n := range foundNames {
			var foundIDs sort.StringSlice = make([]string, len(foundRecords[n]))
			i = 0
			for k := range foundRecords[n] {
				foundIDs[i] = k.String()
				i++
			}

			if asc {
				sort.Sort(foundIDs)
			} else {
				sort.Sort(sort.Reverse(foundIDs))
			}

			for _, id := range foundIDs {
				records = append(records, foundRecords[n][uuid.MustParse(id)])
			}

		}

	default:
		return nil, db.ErrUnknownAttribute
	}

	return &records, nil
}
