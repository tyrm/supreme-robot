package postgres

import (
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/tyrm/supreme-robot/config"
	"reflect"
	"testing"
)

const postgresPort = 23485
const postgresConnection = "host=localhost port=23485 user=postgres password=postgres dbname=postgres sslmode=disable"

func TestNewClient(t *testing.T) {
	database := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().Port(postgresPort))
	if err := database.Start(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := database.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	cfg := config.Config{
		PostgresDsn: postgresConnection,
	}

	client, err := NewClient(&cfg)
	if err != nil {
		t.Fatalf("unexpected error, got: %s, want: nil.", err.Error())
	}
	if reflect.TypeOf(client) != reflect.TypeOf(&Client{}) {
		t.Fatalf("unexpected client type, got: %s, want: %s", reflect.TypeOf(client), reflect.TypeOf(&Client{}))
	}
}
