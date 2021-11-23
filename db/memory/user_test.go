package memory

import (
	"github.com/tyrm/supreme-robot/db/tests"
	"testing"
)

func TestClient_ReadUser_Admin(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	tests.DoReadUserAdmin(t, client)
}

func TestClient_ReadUser_UnknownUser(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	tests.DoReadUserUnknownUser(t, client)
}

func TestClient_ReadUserByUsername_Admin(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	tests.DoReadUserByUsernameAdmin(t, client)
}

func TestClient_ReadUserByUsername_UnknownUser(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	tests.DoReadUserByUsernameUnknownUser(t, client)
}
