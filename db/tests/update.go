package tests

import (
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/models"
	"testing"
)

// DoUpdateUser tests the Update function for a models.User type
func DoUpdateUser(t *testing.T, client db.DB) {
	newUser := models.User{
		Username: "doupdateuser",
	}
	err := newUser.SetPassword("lepassword")
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}
	err = client.Create(&newUser)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	newUser2, err := client.ReadUser(newUser.ID)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	err = newUser2.SetPassword("newpassword")
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	err = client.Update(newUser2)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	newUser3, err := client.ReadUser(newUser.ID)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	if !newUser3.CheckPasswordHash("newpassword") {
		t.Fatalf("password not set properly, tried: 'newpassword'")
	}
}

// DoUpdateUnknownType tests Update for an unknown type
func DoUpdateUnknownType(t *testing.T, client db.DB) {
	newUnknown := unknownType{}
	err := client.Update(&newUnknown)
	if err.Error() != db.ErrUnknownType.Error() {
		t.Fatalf("unexpected error, got: %s, want: %s", err.Error(), db.ErrUnknownType.Error())
	}
}
