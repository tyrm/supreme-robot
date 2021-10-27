package models

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type Config struct {
	Key   string `db:"key"`
	Value string `db:"value"`

	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Model Functions

func (cnf *Config) Create(c *Client) error {
	if cnf.ID != uuid.Nil {
		return ErrAlreadyCreated
	}

	err := c.db.
		QueryRowx(`INSERT INTO "public"."config"("key", "value")
			VALUES ($1, $2) RETURNING id, created_at, updated_at;`, cnf.Key, cnf.Value).
		Scan(&cnf.ID, &cnf.CreatedAt, &cnf.UpdatedAt)

	return err
}

func (cnf *Config) Update(c *Client) error {
	if cnf.ID == uuid.Nil {
		return ErrNotCreated
	}

	err := c.db.
		QueryRowx(`UPDATE public.config SET key=$1, "value"=$2, updated_at=CURRENT_TIMESTAMP
		WHERE id=$3 RETURNING updated_at;`, cnf.Key, cnf.Value, cnf.ID).
		Scan(&cnf.UpdatedAt)

	return err
}

// Client Functions

func (c *Client) ReadConfigByKey(k string) (*Config, error) {
	var config Config
	err := c.db.
		Get(&config, `SELECT id, key, "value", created_at, updated_at
		FROM public.config WHERE key = $1;`, k)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Client) ReadConfigsByKeys(k *[]string) (*[]Config, error) {
	var config []Config
	err := c.db.
		Select(&config, `SELECT id, key, "value", created_at, updated_at
		FROM public.config WHERE key = ANY($1);`, pq.Array(k))
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &config, nil
}
