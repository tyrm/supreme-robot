package memory

import (
	"github.com/tyrm/supreme-robot/db/tests"
	"testing"
)

func TestClient_Delete_Domain(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	tests.DoDeleteDomain(t, client)
}

func TestClient_Delete_UnknownType(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	tests.DoDeleteUnknownType(t, client)
}
