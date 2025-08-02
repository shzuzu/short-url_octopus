package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	const op = "database.db.NewDB"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	err = createTable(db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)

	}

	return &DB{db: db}, nil

}

func createTable(db *sql.DB) error {
	const op = "database.db.CreateTable"
	queries := []string{
		`CREATE TABLE IF NOT EXISTS urls(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL
	);`,
		`CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
func (db *DB) SaveURL(url string, alias string) (int64, error) {
	const op = "database.db.SaveURL"
	stmt, err := db.db.Prepare("INSERT INTO urls(url,alias) values(?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(url, alias)

	if err != nil {
		sqliteErr, ok := err.(sqlite3.Error)
		if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, ErrURLExists)
		}
		return 0, fmt.Errorf("%s: execute op: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get id: %w", op, err)
	}
	return id, nil
}

func (db *DB) GetURL(alias string) (string, error) {
	const op = "database.db.GetURL"

	stmt, err := db.db.Prepare("SELECT url FROM urls WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare get url op: %w", op, err)
	}
	defer stmt.Close()

	var resUrl string
	err = stmt.QueryRow(alias).Scan(&resUrl)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "", ErrURLNotFound
	case err != nil:
		return "", fmt.Errorf("%s: execute op: %w", op, err)

	}
	return resUrl, nil
}

func (db *DB) DeleteURL(alias string) error {
	const op = "database.db.DeleteURL"

	stmt, err := db.db.Prepare("DELETE FROM urls WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare op: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: execute op: %w", op, err)
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: affect rows op: %w", op, err)
	}
	if rowAffected == 0 {
		return ErrURLNotFound
	}
	return nil

}
