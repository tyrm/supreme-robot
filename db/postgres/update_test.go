//go:build integration

package postgres

import (
	"github.com/tyrm/supreme-robot/db/tests"
	"testing"
)

func TestClient_Update_User(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	tests.DoUpdateUser(t, client)
}

func TestClient_Update_UnknownType(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	tests.DoUpdateUnknownType(t, client)
}
