package storage

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"url-shortener/internal/model"
)

type PostgresStorage struct {
	db *sql.DB
}

const createLinksTableSQL = `
CREATE TABLE IF NOT EXISTS links (
    alias        TEXT PRIMARY KEY,
    original_url TEXT NOT NULL,
    clicks       INTEGER NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
)`

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if _, err := db.Exec(createLinksTableSQL); err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Save(link model.Link) error {
	_, err := s.db.Exec(
		`INSERT INTO links (alias, original_url) VALUES ($1, $2)`,
		link.ID, link.URL,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (s *PostgresStorage) FindByID(id string) (model.Link, error) {
	var link model.Link
	err := s.db.QueryRow(
		`SELECT alias, original_url, clicks, created_at FROM links WHERE alias = $1`,
		id,
	).Scan(&link.ID, &link.URL, &link.Clicks, &link.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Link{}, ErrNotFound
	}
	if err != nil {
		return model.Link{}, err
	}
	return link, nil
}

func (s *PostgresStorage) IncrementClicks(id string) error {
	res, err := s.db.Exec(
		`UPDATE links SET clicks = clicks + 1 WHERE alias = $1`,
		id,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
