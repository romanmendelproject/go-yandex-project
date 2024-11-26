package storage

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5"
)

func (pg *PostgresStorage) Register(ctx context.Context, login, password string) (int, error) {
	query := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id;`
	var id int

	if err := pg.db.QueryRow(ctx, query, login, password).Scan(&id); err != nil {
		log.Error("error inserting user", "login", login, "error", err)

		return 0, err
	}

	return id, nil
}

func (pg *PostgresStorage) Login(ctx context.Context, login, password string) (int, error) {
	query := `SELECT id FROM users WHERE login = $1 AND password = $2;`
	var id int

	if err := pg.db.QueryRow(ctx, query, login, password).Scan(&id); err != nil {
		log.Error("error getting user", "login", login, "error", err)

		return 0, err
	}

	return id, nil
}

func (pg *PostgresStorage) CheckLogin(ctx context.Context, login string) error {
	var userID int

	if err := pg.db.QueryRow(ctx, `SELECT id FROM users WHERE login=$1`, login).Scan(&userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}

		log.Error("error querying user", " login ", login, " error ", err)

		return err
	}

	if userID != 0 {
		log.Error("user already exists", "login", login, "error")

		return errors.New("user already exists")
	}

	return nil
}
