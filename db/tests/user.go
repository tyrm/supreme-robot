package tests

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"reflect"
	"testing"
)

var userAdmin = models.User{
	ID:       uuid.MustParse("8c504483-1e11-4243-b6c8-14499877a641"),
	Username: "admin",
	Password: "password",
}

// DoReadUserAdmin tests the ReadUser function for a known user
func DoReadUserAdmin(t *testing.T, client db.DB) {
	receivedUser, err := client.ReadUser(userAdmin.ID)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil", err.Error())
	}

	if receivedUser.ID != userAdmin.ID {
		t.Errorf("unexpected username, got: %s, want: %s", receivedUser.ID, userAdmin.ID)
	}

	if receivedUser.Username != userAdmin.Username {
		t.Errorf("unexpected username, got: %s, want: %s", receivedUser.Username, userAdmin.Username)
	}

	if !receivedUser.CheckPasswordHash(userAdmin.Password) {
		t.Errorf("invalid password, tried: '%s'", userAdmin.Password)
	}

	if !receivedUser.IsMemberOfGroup([]uuid.UUID{models.GroupSuperAdmin}...) {
		t.Error("user missing superadmin group")
	}
}

// DoReadUserUnknownUser tests the ReadUser function for an unknown user
func DoReadUserUnknownUser(t *testing.T, client db.DB) {
	id := uuid.MustParse("5cacb040-fa25-4149-9829-e184b4dbcbde")
	receivedUser, err := client.ReadUser(id)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil", err.Error())
	}
	if receivedUser != nil {
		t.Fatalf("unexpected user, got: %s, want: nil", err.Error())
	}
}

// DoReadUserByUsernameAdmin tests the ReadUserByUsername function for a known user
func DoReadUserByUsernameAdmin(t *testing.T, client db.DB) {
	receivedUser, err := client.ReadUserByUsername(userAdmin.Username)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}
	if reflect.TypeOf(receivedUser) != reflect.TypeOf(&models.User{}) {
		t.Errorf("unexpected type, got: %s, want: %s", reflect.TypeOf(receivedUser), reflect.TypeOf(&models.User{}))
		return
	}
	if receivedUser == nil {
		t.Errorf("expected object, got: nil")
		return
	}

	if receivedUser.ID != userAdmin.ID {
		t.Errorf("unexpected username, got: %s, want: %s", receivedUser.ID, userAdmin.ID)
	}

	if receivedUser.Username != userAdmin.Username {
		t.Errorf("unexpected username, got: %s, want: %s", receivedUser.Username, userAdmin.Username)
	}

	if !receivedUser.CheckPasswordHash(userAdmin.Password) {
		t.Errorf("invalid password, tried: '%s'", userAdmin.Password)
	}

	if !receivedUser.IsMemberOfGroup([]uuid.UUID{models.GroupSuperAdmin}...) {
		t.Error("user missing superadmin group")
	}
}

// DoReadUserByUsernameUnknownUser tests the ReadUserByUsername function for an unknown user
func DoReadUserByUsernameUnknownUser(t *testing.T, client db.DB) {
	receivedUser, err := client.ReadUserByUsername("invaliduser")
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil", err.Error())
	}
	if receivedUser != nil {
		t.Fatalf("unexpected user, got: %s, want: nil", err.Error())
	}
}
