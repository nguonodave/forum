package models

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Database connection and initialization.
// go : embed migrations .sql
var fs embed.FS

const (
	dbName = "forum.db"
)

func InitializeDB() (*sql.DB, error) {
	// open database connection
	// data source name
	dsn := fmt.Sprintf("file:%s?_foreign_keys=on&_journal_mode=WAL", dbName)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to open databas: %w", err)
	}

	// verify connection to database
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}
	// apply  migrations
	if err := applyMigration(db); err != nil {
		return nil, fmt.Errorf("Failed to apply migrations: %w", err)
	}
	// configure connection pool
	db.SetMaxOpenConns(1) // sqlite only supports 1 wrter  at a time
	db.SetMaxIdleConns(1)

	return db, nil
}

// Apply migrations to the database.
func applyMigration(db *sql.DB) error {
}
