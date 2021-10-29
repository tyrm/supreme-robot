package models

import (
	"github.com/jmoiron/sqlx"
	"github.com/juju/loggo"
	_ "github.com/lib/pq"
	"github.com/markbates/pkger"
	"github.com/rubenv/sql-migrate"
	"github.com/tyrm/supreme-robot/startup"
)

var logger = loggo.GetLogger("models")

type Client struct {
	db *sqlx.DB
}

func NewClient(cfg *startup.Config) (*Client, error) {
	client, err := sqlx.Connect("postgres", cfg.PostgresDsn)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Loading Migrations")
	migrations := &migrate.HttpFileSystemMigrationSource{
		FileSystem: pkger.Dir("/models/migrations"),
	}

	logger.Debugf("Applying Migrations")
	n, err := migrate.Exec(client.DB, "postgres", migrations, migrate.Up)
	if n > 0 {
		logger.Infof("Applied %d migrations!", n)
	}
	if err != nil {
		logger.Criticalf("Could not migrate models: %s", err)
		return nil, err
	}

	return &Client{
		db: client,
	}, nil
}
