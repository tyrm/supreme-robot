package memory

import (
	"github.com/tyrm/supreme-robot/kv/tests"
	"testing"
)

func TestClient_AddDomain(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	tests.DoAddDomain(t, client)
}

func TestClient_GetDomains(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	tests.DoGetDomain(t, client)
}

func TestClient_RemoveDomain(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	tests.DoRemoveDomain(t, client)
}
