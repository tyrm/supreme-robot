package memory

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/db"
	"testing"
)

func TestClient_Update_User(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	adminUser, err := client.ReadUser(uuid.MustParse("44892097-2c97-4c16-b4d1-e8522586df48"))
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

	adminUser, err = client.ReadUser(uuid.MustParse("44892097-2c97-4c16-b4d1-e8522586df48"))
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	if !adminUser.CheckPasswordHash("newpassword") {
		t.Fatalf("password not set properly, tried: 'newpassword'")
	}
}

func TestClient_Update_UnknownType(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}

	newUnknown := unknownType{}
	err = client.Update(&newUnknown)
	if err.Error() != db.ErrUnknownType.Error() {
		t.Fatalf("unexpected error, got: %s, want: %s", err.Error(), db.ErrUnknownType.Error())
	}
}
