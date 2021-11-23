//go:build integration

package postgres

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/config"
	"github.com/tyrm/supreme-robot/models"
	"os"
	"reflect"
	"testing"
)

func TestClient_ReadUser_Admin(t *testing.T) {
	cfg := config.Config{
		PostgresDsn: os.Getenv("TEST_DSN"),
	}
	client, err := NewClient(&cfg)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	id := uuid.MustParse("8c504483-1e11-4243-b6c8-14499877a641")
	receivedUser, err := client.ReadUser(id)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}
	if reflect.TypeOf(receivedUser) != reflect.TypeOf(&models.User{}) {
		t.Errorf("unexpected type, got: %s, want: %s", reflect.TypeOf(receivedUser), reflect.TypeOf(&models.User{}))
		return
	}
	if receivedUser == nil {
		t.Errorf("expected object, got: nil")
		return
	}

	if receivedUser.ID != uuid.MustParse("8c504483-1e11-4243-b6c8-14499877a641") {
		t.Errorf("unexpected username, got: %s, want: 8c504483-1e11-4243-b6c8-14499877a641", receivedUser.ID)
	}

	if receivedUser.Username != "admin" {
		t.Errorf("unexpected username, got: %s, want: admin", receivedUser.Username)
	}

	if !receivedUser.CheckPasswordHash("password") {
		t.Error("invalid password, tried: 'password'")
	}

	if !receivedUser.IsMemberOfGroup([]uuid.UUID{models.GroupSuperAdmin}...) {
		t.Error("user missing superadmin group")
	}
}

func TestClient_ReadUser_UnknownUser(t *testing.T) {
	cfg := config.Config{
		PostgresDsn: os.Getenv("TEST_DSN"),
	}
	client, err := NewClient(&cfg)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
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
	cfg := config.Config{
		PostgresDsn: os.Getenv("TEST_DSN"),
	}
	client, err := NewClient(&cfg)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	receivedUser, err := client.ReadUserByUsername("admin")
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

	if receivedUser.ID != uuid.MustParse("8c504483-1e11-4243-b6c8-14499877a641") {
		t.Errorf("unexpected username, got: %s, want: 8c504483-1e11-4243-b6c8-14499877a641", receivedUser.ID)
	}

	if receivedUser.Username != "admin" {
		t.Errorf("unexpected username, got: %s, want: admin", receivedUser.Username)
	}

	if !receivedUser.CheckPasswordHash("password") {
		t.Error("invalid password, tried: 'password'")
	}

	if !receivedUser.IsMemberOfGroup([]uuid.UUID{models.GroupSuperAdmin}...) {
		t.Error("user missing superadmin group")
	}
}

func TestClient_ReadUserByUsername_UnknownUser(t *testing.T) {
	cfg := config.Config{
		PostgresDsn: os.Getenv("TEST_DSN"),
	}
	client, err := NewClient(&cfg)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	receivedUser, err := client.ReadUserByUsername("invaliduser")
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil", err.Error())
	}
	if receivedUser != nil {
		t.Fatalf("unexpected user, got: %s, want: nil", err.Error())
	}
}
