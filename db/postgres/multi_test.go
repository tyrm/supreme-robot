//go:build integration

package postgres

import (
	"github.com/tyrm/supreme-robot/db/tests"
	"testing"
)

func TestClient_CreateDomainWRecords(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	tests.DoCreateDomainWRecords(t, client)
}

func TestClient_CreateGroupsForUser(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	tests.DoCreateGroupsForUser(t, client)
}

func TestClient_CreateGroupsForUser_UnknownUser(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	tests.DoCreateGroupsForUserUnknownUser(t, client)
}
