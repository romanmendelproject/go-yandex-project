package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"
)

func init() {
	goose.AddMigrationContext(upInitTable, downInitTable)
}

func upInitTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	log.Info("Create DB Table")

	query := `
		CREATE TABLE IF NOT EXISTS users (
		    id SERIAL PRIMARY KEY,
		    login varchar(255) NOT NULL,
		    password varchar(255) NOT NULL,
		    current FLOAT DEFAULT 0,
		    withdrawn FLOAT DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS cred (
			name CHAR(30),
			username CHAR(30),
			password BYTEA,
			meta CHAR(30),
			user_id integer REFERENCES users
		);

		CREATE TABLE IF NOT EXISTS text (
			name CHAR(30),
			data CHAR(500),
			meta CHAR(30),
			user_id integer REFERENCES users
		);

		CREATE TABLE IF NOT EXISTS byte (
			name CHAR(30),
			data BYTEA,
			meta CHAR(30),
			user_id integer REFERENCES users
		);

		CREATE TABLE IF NOT EXISTS card (
			name CHAR(30),
			data BIGINT,
			meta CHAR(30),
			user_id integer REFERENCES users
		);
	`

	// Creating metrics table
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downInitTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	log.Info("Remove DB Table")

	if _, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS packets"); err != nil {
		return err
	}

	return nil
}
