package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func NewInMemory() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "file:users?mode=memory&cache=shared")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	statement := `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	role TEXT NOT NULL
);`

	if _, err := db.Exec(statement); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	return nil
}

func SeedUsers(db *sql.DB) error {
	statement := `
INSERT INTO users (id, name, email, role) VALUES
	(1, 'Ada Lovelace', 'ada@example.com', 'admin'),
	(2, 'Grace Hopper', 'grace@example.com', 'engineer'),
	(3, 'Alan Turing', 'alan@example.com', 'researcher')
ON CONFLICT(id) DO NOTHING;`

	if _, err := db.Exec(statement); err != nil {
		return fmt.Errorf("seed users: %w", err)
	}

	return nil
}
