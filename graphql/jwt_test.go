package graphql

import (
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/models"
	"testing"
)

func TestCreateToken(t *testing.T) {
	// create server
	server, _, _, _, err := newTestServer()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}
	if server == nil {
		t.Errorf("expected server, got: nil, want: *Server.")
	}

	// Create User
	user := models.User{
		ID: uuid.Must(uuid.Parse("faf3da6e-afab-4e49-99b6-483da3aa4226")),
	}

	// Create Token
	token, err := server.createToken(&user)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}

	// verify data
	if token.AccessUUID == uuid.Nil {
		t.Errorf("invalid access uuid, got: %s, want: something else.", token.AccessUUID)
	}

	rID := token.AccessUUID.String() + "++" + user.ID.String()
	if token.RefreshUUID != rID {
		t.Errorf("invalid access uuid, got: %s, want: %s.", token.AccessUUID, rID)
	}

	t.Logf("%#v", token)
}
