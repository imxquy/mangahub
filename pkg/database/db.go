package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE,
			email TEXT UNIQUE,
			password_hash TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS manga (
			id TEXT PRIMARY KEY,
			title TEXT,
			author TEXT,
			genres TEXT,
			status TEXT,
			total_chapters INTEGER,
			description TEXT,
			cover_url TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS user_progress (
			user_id TEXT,
			manga_id TEXT,
			current_chapter INTEGER,
			status TEXT,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, manga_id)
		);`,
	}
	for i, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("migrate stmt %d: %w", i, err)
		}
	}
	return nil
}
