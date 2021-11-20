package memory

import (
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/tyrm/supreme-robot/kv"
	"testing"
	"time"
)

func TestClient_DeleteAccessToken(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	accessToken := uuid.MustParse("27efbdd8-0ecc-4922-a206-ea9d5a71a1a6")
	userID := uuid.MustParse("1e82d82e-29fc-41f8-9c8d-fd495cf2fea2")

	client.KV.Set(kv.KeyJwtAccess(accessToken.String()), userID.String(), cache.NoExpiration)

	_, err = client.DeleteAccessToken(accessToken)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	result, resultFound := client.KV.Get(kv.KeyJwtAccess(accessToken.String()))
	if resultFound {
		t.Errorf("delete failed, got: %v, want: nil.", result)
	}
}

func TestClient_DeleteRefreshToken(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	accessToken := uuid.MustParse("27efbdd8-0ecc-4922-a206-ea9d5a71a1a6")
	userID := uuid.MustParse("1e82d82e-29fc-41f8-9c8d-fd495cf2fea2")

	refreshToken := accessToken.String() + "++" + userID.String()

	client.KV.Set(kv.KeyJwtRefresh(refreshToken), userID.String(), cache.NoExpiration)

	_, err = client.DeleteRefreshToken(refreshToken)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	result, resultFound := client.KV.Get(kv.KeyJwtRefresh(refreshToken))
	if resultFound {
		t.Errorf("delete failed, got: %v, want: nil.", result)
	}
}

func TestClient_GetAccessToken(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	accessToken := uuid.MustParse("27efbdd8-0ecc-4922-a206-ea9d5a71a1a6")
	userID := uuid.MustParse("1e82d82e-29fc-41f8-9c8d-fd495cf2fea2")

	client.KV.Set(kv.KeyJwtAccess(accessToken.String()), userID.String(), cache.NoExpiration)

	receivedUserID, err := client.GetAccessToken(accessToken)
	if receivedUserID != userID {
		t.Errorf("unexpected user id, got: %v, want: %v", receivedUserID, userID)
	}
}

func TestClient_SetAccessToken(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	accessToken := uuid.MustParse("27efbdd8-0ecc-4922-a206-ea9d5a71a1a6")
	userID := uuid.MustParse("1e82d82e-29fc-41f8-9c8d-fd495cf2fea2")

	err = client.SetAccessToken(accessToken, userID, 10*time.Second)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	receivedUserID, resultFound := client.KV.Get(kv.KeyJwtAccess(accessToken.String()))
	if !resultFound {
		t.Errorf("set failed, got: nil, want: %s.", userID.String())
		return
	}
	receivedUserUUID, err := uuid.Parse(receivedUserID.(string))
	if err != nil {
		t.Errorf("error parsing id %s, got: %s, want: nil.", receivedUserID, err.Error())
		return
	}
	if receivedUserUUID.String() != userID.String() {
		t.Errorf("unexpected user id, got: %v, want: %v", receivedUserUUID, userID)
	}
}

func TestClient_SetRefreshToken(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	accessToken := uuid.MustParse("27efbdd8-0ecc-4922-a206-ea9d5a71a1a6")
	userID := uuid.MustParse("1e82d82e-29fc-41f8-9c8d-fd495cf2fea2")

	refreshToken := accessToken.String() + "++" + userID.String()

	err = client.SetRefreshToken(refreshToken, userID, 10*time.Second)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	receivedUserID, resultFound := client.KV.Get(kv.KeyJwtRefresh(refreshToken))
	if !resultFound {
		t.Errorf("set failed, got: nil, want: %s.", userID.String())
		return
	}
	receivedUserUUID, err := uuid.Parse(receivedUserID.(string))
	if err != nil {
		t.Errorf("error parsing id %s, got: %s, want: nil.", receivedUserID, err.Error())
		return
	}
	if receivedUserUUID.String() != userID.String() {
		t.Errorf("unexpected user id, got: %v, want: %v", receivedUserUUID, userID)
	}
}
