//go:build integration

package postgres

import (
	"github.com/tyrm/supreme-robot/config"
	"os"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := testCreateClient()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
	}
	if reflect.TypeOf(client) != reflect.TypeOf(&Client{}) {
		t.Errorf("unexpected client type, got: %s, want: %s", reflect.TypeOf(client), reflect.TypeOf(&Client{}))
	}
}

func testCreateClient() (*Client, error) {
	cfg := config.Config{
		PostgresDsn: os.Getenv("TEST_DSN"),
	}

	return NewClient(&cfg)
}
