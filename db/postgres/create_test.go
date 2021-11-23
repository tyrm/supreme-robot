//go:build integration

package postgres

import (
	"github.com/tyrm/supreme-robot/db/tests"
	"testing"
)

func TestClient_Create_Domain(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	tests.DoCreateDomain(t, client)
}

func TestClient_Create_Record(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	tests.DoCreateRecord(t, client)
}

func TestClient_Create_User(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	tests.DoCreateUser(t, client)
}

func TestClient_Create_UnknownType(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	tests.DoCreateUnknownType(t, client)
}
