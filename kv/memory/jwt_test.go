package memory

import (
	"github.com/tyrm/supreme-robot/kv/tests"
	"testing"
)

func TestClient_DeleteAccessToken(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	tests.DoDeleteAccessToken(t, client)
}

func TestClient_DeleteRefreshToken(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	tests.DoDeleteRefreshToken(t, client)
}

func TestClient_GetAccessToken(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	tests.DoGetAccessToken(t, client)
}

func TestClient_GetAccessToken_NotFound(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	tests.DoGetAccessTokenNotFound(t, client)
}

func TestClient_SetAccessToken(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	tests.DoSetAccessToken(t, client)
}

func TestClient_SetRefreshToken(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	tests.DoSetRefreshToken(t, client)
}
