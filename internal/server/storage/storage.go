package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/AlexanderGrom/componenta/crypt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romanmendelproject/go-yandex-project/internal/types"
	log "github.com/sirupsen/logrus"
)

type PostgresStorage struct {
	db *pgxpool.Pool
}

var (
	pgInstance *PostgresStorage
	pgOnce     sync.Once
)

func NewPostgresStorage(ctx context.Context, connString string) *PostgresStorage {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			log.Fatal("unable to create connection pool: %w", err)
		}

		pgInstance = &PostgresStorage{db}
	})

	return pgInstance
}

func (pg *PostgresStorage) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *PostgresStorage) Close() {
	pg.db.Close()
}

// GetCred читает данные формата CredType из БД
func (pg *PostgresStorage) GetCred(ctx context.Context, name string, userID int) (*types.CredType, error) {
	var (
		values   types.CredType
		username sql.NullString
		password sql.NullString
		meta     string
	)

	if err := pg.db.QueryRow(ctx, "SELECT username, password, meta FROM cred WHERE name = $1 AND user_id = $2 ", name, userID).Scan(&username, &password, &meta); err != nil {
		return nil, err
	}

	if !username.Valid || !password.Valid {
		return nil, fmt.Errorf("unexpected type of cred")
	}

	passwordString, err := crypt.Decrypt(password.String, "Secret_Key")

	if err != nil {
		log.Fatalln("Decrypt:", err)
	}

	values.Username = username.String
	values.Password = passwordString
	values.Meta = meta

	return &values, nil
}

// SetCred записывает данные формата CredType в БД
func (pg *PostgresStorage) SetCred(ctx context.Context, value types.CredType, userID int) error {
	var (
		username sql.NullString
		password sql.NullString
		meta     string
	)

	tx, err := pg.db.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	passwordHash, err := crypt.Encrypt(value.Password, "Secret_Key")

	if err != nil {
		log.Fatalln("Encrypt:", err)
	}

	// Check if data exists
	if err := tx.QueryRow(ctx, "SELECT username, password, meta FROM cred WHERE name = $1 AND user_id = $2", value.Name, userID).Scan(&username, &password, &meta); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Insert new metric if not exists
			log.Error(value.Name)
			if _, err := tx.Exec(ctx, `INSERT INTO cred (name, username, password, meta, user_id) VALUES ($1, $2, $3, $4, $5)`, value.Name, value.Username, passwordHash, value.Meta, userID); err != nil {
				log.Error(err)
				return err
			}
			return nil
		}

		log.Error(err)
		return err
	}

	// Update data if exists
	if _, err := tx.Exec(ctx, `UPDATE cred SET username = $2, password = $3, meta = $4 WHERE name = $1`, value.Name, value.Username, passwordHash, value.Meta); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// GetText читает данные формата TextType из БД
func (pg *PostgresStorage) GetText(ctx context.Context, name string, userID int) (*types.TextType, error) {
	var (
		values types.TextType
		data   sql.NullString
		meta   string
	)

	if err := pg.db.QueryRow(ctx, "SELECT data, meta FROM text WHERE name = $1 AND user_id = $2 ", name, userID).Scan(&data, &meta); err != nil {
		return nil, err
	}

	if !data.Valid {
		return nil, fmt.Errorf("unexpected type of text")
	}

	values.Data = data.String
	values.Meta = meta

	return &values, nil
}

// SetText записывает данные формата TextType в БД
func (pg *PostgresStorage) SetText(ctx context.Context, value types.TextType, userID int) error {
	var (
		data sql.NullString
		meta string
	)
	tx, err := pg.db.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	// Check if data exists
	if err := tx.QueryRow(ctx, "SELECT data, meta FROM text WHERE name = $1 AND user_id = $2", value.Name, userID).Scan(&data, &meta); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Insert new metric if not exists
			log.Error(value.Name)
			if _, err := tx.Exec(ctx, `INSERT INTO text (name, data, meta, user_id) VALUES ($1, $2, $3, $4)`, value.Name, value.Data, value.Meta, userID); err != nil {
				log.Error(err)
				return err
			}
			return nil
		}

		log.Error(err)
		return err
	}

	// Update data if exists
	if _, err := tx.Exec(ctx, `UPDATE text SET data = $2, meta = $3 WHERE name = $1`, value.Name, value.Data, value.Meta); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// GetByte читает данные формата ByteType из БД
func (pg *PostgresStorage) GetByte(ctx context.Context, name string, userID int) (*types.ByteType, error) {
	var (
		values types.ByteType
		data   sql.RawBytes
		meta   string
	)

	if err := pg.db.QueryRow(ctx, "SELECT data, meta FROM byte WHERE name = $1 AND user_id = $2 ", name, userID).Scan(&data, &meta); err != nil {
		return nil, err
	}

	values.Data = data
	values.Meta = meta

	return &values, nil
}

// SetTByte записывает данные формата TextType в БД
func (pg *PostgresStorage) SetByte(ctx context.Context, value types.ByteType, userID int) error {
	var (
		data sql.RawBytes
		meta string
	)

	tx, err := pg.db.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	// Check if data exists
	if err := tx.QueryRow(ctx, "SELECT data, meta FROM byte WHERE name = $1 AND user_id = $2", value.Name, userID).Scan(&data, &meta); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Insert new metric if not exists
			log.Error(value.Name)
			if _, err := tx.Exec(ctx, `INSERT INTO byte (name, data, meta, user_id) VALUES ($1, $2, $3, $4)`, value.Name, value.Data, value.Meta, userID); err != nil {
				log.Error(err)
				return err
			}
			return nil
		}

		log.Error(err)
		return err
	}

	// Update data if exists
	if _, err := tx.Exec(ctx, `UPDATE byte SET data = $2, meta = $3 WHERE name = $1`, value.Name, value.Data, value.Meta); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// GetCard читает данные формата CardType из БД
func (pg *PostgresStorage) GetCard(ctx context.Context, name string, userID int) (*types.CardType, error) {
	var (
		values types.CardType
		data   sql.NullString
		meta   string
	)

	if err := pg.db.QueryRow(ctx, "SELECT data, meta FROM card WHERE name = $1 AND user_id = $2 ", name, userID).Scan(&data, &meta); err != nil {
		return nil, err
	}

	if !data.Valid {
		return nil, fmt.Errorf("unexpected type of cred")
	}

	dataString, err := crypt.Decrypt(data.String, "Secret_Key")

	if err != nil {
		log.Fatalln("Decrypt:", err)
	}

	dataInt, err := strconv.ParseInt(dataString, 10, 64)
	if err != nil {
		panic(err)
	}

	values.Data = dataInt
	values.Meta = meta

	return &values, nil
}

// SetCard записывает данные формата CardType в БД
func (pg *PostgresStorage) SetCard(ctx context.Context, value types.CardType, userID int) error {
	var (
		data sql.NullInt64
		meta string
	)

	tx, err := pg.db.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	dataHash, err := crypt.Encrypt(strconv.FormatInt(value.Data, 10), "Secret_Key")

	if err != nil {
		log.Fatalln("Encrypt:", err)
	}

	// Check if data exists
	if err := tx.QueryRow(ctx, "SELECT data, meta FROM card WHERE name = $1 AND user_id = $2", value.Name, userID).Scan(&data, &meta); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Insert new metric if not exists
			log.Error(value.Name)
			if _, err := tx.Exec(ctx, `INSERT INTO card (name, data, meta, user_id) VALUES ($1, $2, $3, $4)`, value.Name, dataHash, value.Meta, userID); err != nil {
				log.Error(err)
				return err
			}
			return nil
		}

		log.Error(err)
		return err
	}

	// Update data if exists
	if _, err := tx.Exec(ctx, `UPDATE card SET data = $2, meta = $3 WHERE name = $1`, value.Name, dataHash, value.Meta); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
