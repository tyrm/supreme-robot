package memory

import (
	"github.com/tyrm/supreme-robot/db/tests"
	"testing"
)

func TestClient_ReadRecordsForDomain(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	tests.DoReadRecordsForDomain(t, client)
}

func TestClient_ReadRecordsForDomain_OrderBy_Unknown(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	tests.DoReadRecordsForDomainOrderByUnknown(t, client)
}

func TestClient_ReadRecordsForDomainByName(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	tests.DoReadRecordsForDomainByName(t, client)
}
