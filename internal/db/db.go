package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS logs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	task_name TEXT NOT NULL,
	estimate_hours REAL NOT NULL,
	actual_hours REAL NOT NULL,
	ai_minutes INTEGER,
	problem TEXT,
	solution TEXT,
	learning TEXT
);
`

// DB wraps the SQL database connection
type DB struct {
	*sql.DB
}

// Open opens or creates the database at ~/.djou/djou.db
func Open() (*DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dbDir := filepath.Join(homeDir, ".djou")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dbDir, "djou.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}

	return &DB{db}, nil
}
