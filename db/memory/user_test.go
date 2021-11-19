package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/models"
	"testing"
)

func TestClient_ReadUser_Admin(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	id := uuid.MustParse("44892097-2c97-4c16-b4d1-e8522586df48")
	receivedUser, err := client.ReadUser(id)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	if receivedUser.ID != uuid.MustParse("44892097-2c97-4c16-b4d1-e8522586df48") {
		t.Errorf("unexpected username, got: %s, want: 44892097-2c97-4c16-b4d1-e8522586df48", err.Error())
	}

	if receivedUser.Username != "admin" {
		t.Errorf("unexpected username, got: %s, want: admin", err.Error())
	}

	if !receivedUser.CheckPasswordHash("password") {
		t.Error("invalid password, tried: 'password'")
	}

	if !receivedUser.IsMemberOfGroup([]uuid.UUID{models.GroupSuperAdmin}...) {
		t.Error("user missing superadmin group")
	}
}

func TestClient_ReadUser_UnknownUser(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	id := uuid.MustParse("5cacb040-fa25-4149-9829-e184b4dbcbde")
	receivedUser, err := client.ReadUser(id)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil", err.Error())
	}
	if receivedUser != nil {
		t.Fatalf("unexpected user, got: %s, want: nil", err.Error())
	}
}

func TestClient_ReadUserByUsername_Admin(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	receivedUser, err := client.ReadUserByUsername("admin")
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	if receivedUser.ID != uuid.MustParse("44892097-2c97-4c16-b4d1-e8522586df48") {
		t.Errorf("unexpected username, got: %s, want: 44892097-2c97-4c16-b4d1-e8522586df48", err.Error())
	}

	if receivedUser.Username != "admin" {
		t.Errorf("unexpected username, got: %s, want: admin", err.Error())
	}

	if !receivedUser.CheckPasswordHash("password") {
		t.Error("invalid password, tried: 'password'")
	}

	if !receivedUser.IsMemberOfGroup([]uuid.UUID{models.GroupSuperAdmin}...) {
		t.Error("user missing superadmin group")
	}
}

func TestClient_ReadUserByUsername_UnknownUser(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	receivedUser, err := client.ReadUserByUsername("invaliduser")
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil", err.Error())
	}
	if receivedUser != nil {
		t.Fatalf("unexpected user, got: %s, want: nil", err.Error())
	}
}
