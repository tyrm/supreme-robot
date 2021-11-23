//go:build integration

package postgres

import (
	"github.com/tyrm/supreme-robot/db/tests"
	"testing"
)

func TestClient_ReadDomain(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	tests.DoReadDomain(t, client)
}

func TestClient_ReadDomain_Unknown(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	tests.DoReadDomainUnknown(t, client)
}

func TestClient_ReadDomainZ(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	tests.DoReadDomainZ(t, client)
}

func TestClient_ReadDomainZ_Unknown(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	tests.DoReadDomainZUnknown(t, client)
}

func TestClient_ReadDomainByDomain(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	tests.DoReadDomainByDomain(t, client)
}

func TestClient_ReadDomainByDomain_Unknown(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil", err.Error())
		return
	}

	tests.DoReadDomainByDomainUnknown(t, client)
}

func TestClient_ReadDomainsForUser(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	tests.DoReadDomainsForUser(t, client)
}
