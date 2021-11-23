package tests

import (
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/tyrm/supreme-robot/kv"
	"testing"
	"time"
)

// DoDeleteAccessToken tests the DeleteAccessToken function
func DoDeleteAccessToken(t *testing.T, client kv.Webapp) {

	accessToken := uuid.MustParse("b2758ab9-e795-4b93-b688-fdf899cc53a5")
	userID := uuid.MustParse("5f1f9747-3a3f-4c29-ad37-2fddc06cb319")

	err := client.SetAccessToken(accessToken, userID, cache.NoExpiration)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	_, err = client.DeleteAccessToken(accessToken)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	result, err := client.GetAccessToken(accessToken)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}
	if result != uuid.Nil {
		t.Errorf("delete failed, got: %v, want: nil.", result)
	}
}

// DoDeleteRefreshToken tests the DeleteRefreshToken function
func DoDeleteRefreshToken(t *testing.T, client kv.Webapp) {
	accessToken := uuid.MustParse("56fb68a3-8ec0-4b67-95f2-88fd2b7f8ebd")
	userID := uuid.MustParse("d5bc6797-119e-4ed6-b029-fcce1e3cae7a")

	refreshToken := accessToken.String() + "++" + userID.String()

	err := client.SetRefreshToken(refreshToken, userID, cache.NoExpiration)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	_, err = client.DeleteRefreshToken(refreshToken)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	result, err := client.GetRefreshToken(refreshToken)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}
	if result != uuid.Nil {
		t.Errorf("delete failed, got: %v, want: nil.", result)
	}
}

// DoGetAccessToken tests the GetAccessToken function
func DoGetAccessToken(t *testing.T, client kv.Webapp) {
	accessToken := uuid.MustParse("9cd471c6-f5d8-4d9d-921c-f7a60c7fd645")
	userID := uuid.MustParse("8b634db6-7764-4698-8a41-bcd1475688e9")

	err := client.SetAccessToken(accessToken, userID, 24*time.Hour)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	receivedUserID, err := client.GetAccessToken(accessToken)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}
	if receivedUserID != userID {
		t.Errorf("unexpected user id, got: %v, want: %v", receivedUserID, userID)
	}
}

// DoGetAccessTokenNotFound tests the GetAccessToken for an unknown token
func DoGetAccessTokenNotFound(t *testing.T, client kv.Webapp) {
	accessToken := uuid.MustParse("880c2531-44b8-4d04-b3b0-bd24aea40168")

	receivedUserID, err := client.GetAccessToken(accessToken)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}
	if receivedUserID != uuid.Nil {
		t.Errorf("unexpected user id, got: %v, want: nil", receivedUserID)
	}
}

// DoSetAccessToken tests the SetAccessToken function
func DoSetAccessToken(t *testing.T, client kv.Webapp) {
	accessToken := uuid.MustParse("eb95d366-d178-4051-80f6-ac5fa6959d0a")
	userID := uuid.MustParse("824bbb19-27b8-4e49-8620-c036bc1bf70e")

	err := client.SetAccessToken(accessToken, userID, 10*time.Second)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	receivedUserID, err := client.GetAccessToken(accessToken)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}
	if receivedUserID != userID {
		t.Errorf("unexpected user id, got: %s, want: %s", receivedUserID, userID)
	}

	time.Sleep(20 * time.Second)
	receivedUserID2, err := client.GetAccessToken(accessToken)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}
	if receivedUserID2 != uuid.Nil {
		t.Errorf("expiration failed, got: %s, want: nil.", receivedUserID2)
		return
	}
}

// DoGetRefreshTokenNotFound tests the GetAccessToken for an unknown token
func DoGetRefreshTokenNotFound(t *testing.T, client kv.Webapp) {
	accessToken := uuid.MustParse("9c90aabe-c420-4c6f-9f9f-63f5a84244c9")
	userID := uuid.MustParse("f1323130-8263-4386-98a6-2b9eb1749e0d")

	refreshToken := accessToken.String() + "++" + userID.String()

	receivedUserID, err := client.GetRefreshToken(refreshToken)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}
	if receivedUserID != uuid.Nil {
		t.Errorf("unexpected user id, got: %v, want: nil", receivedUserID)
	}
}

// DoSetRefreshToken tests the SetRefreshToken function
func DoSetRefreshToken(t *testing.T, client kv.Webapp) {
	accessToken := uuid.MustParse("b1f14b0a-a780-490c-a7e3-0d54814b86d9")
	userID := uuid.MustParse("e39288b7-1b69-4fb3-be4a-d693adc9c2cd")

	refreshToken := accessToken.String() + "++" + userID.String()

	err := client.SetRefreshToken(refreshToken, userID, 10*time.Second)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	receivedUserID, err := client.GetRefreshToken(refreshToken)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}
	if receivedUserID != userID {
		t.Errorf("unexpected user id, got: %s, want: %s", receivedUserID, userID)
	}

	time.Sleep(20 * time.Second)
	receivedUserID2, err := client.GetRefreshToken(refreshToken)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}
	if receivedUserID2 != uuid.Nil {
		t.Errorf("expiration failed, got: %s, want: nil.", receivedUserID2)
		return
	}
}
