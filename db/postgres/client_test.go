//go:build integration

package postgres

import (
	"github.com/tyrm/supreme-robot/config"
	"os"
	"reflect"
	"testing"
)

var testClient *Client

func TestNewClient(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}
	if reflect.TypeOf(client) != reflect.TypeOf(&Client{}) {
		t.Fatalf("unexpected client type, got: %s, want: %s", reflect.TypeOf(client), reflect.TypeOf(&Client{}))
	}
	if client == nil {
		t.Fatalf("expected client, got: nil")
	}
}

func testCreateClient() (*Client, error) {
	if testClient != nil {
		return testClient, nil
	}

	cfg := config.Config{
		PostgresDsn: os.Getenv("TEST_DSN"),
	}
	newClient, err := NewClient(&cfg)
	if err != nil {
		return nil, err
	}
	testClient = newClient

	return testClient, nil
}
