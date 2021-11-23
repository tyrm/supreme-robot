package tests

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/db"
	"testing"
)

// DoUpdateUser tests the Update function for a models.User type
func DoUpdateUser(t *testing.T, client db.DB) {
	adminUser, err := client.ReadUser(userAdmin.ID)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	err = adminUser.SetPassword("newpassword")
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	err = client.Update(adminUser)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	adminUser, err = client.ReadUser(uuid.MustParse("8c504483-1e11-4243-b6c8-14499877a641"))
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	if !adminUser.CheckPasswordHash("newpassword") {
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
