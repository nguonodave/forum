package database

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Database connection and initialization.
//
//go:embed schema.sql
var fs embed.FS

const (
	dbName     = "forum.db"
	schemaFile = "schema.sql"
)

func InitializeDB() (*sql.DB, error) {
	// open database connection
	// data source name
	dsn := fmt.Sprintf("file:%s?_foreign_keys=on&_journal_mode=WAL", dbName)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// verify connection to database
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	// apply  migrations
	if err := applyMigration(db); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}
	// configure connection pool
	db.SetMaxOpenConns(1) // sqlite only supports 1 wrter  at a time
	db.SetMaxIdleConns(1)

	return db, nil
}

// Apply migrations to the database.
func applyMigration(db *sql.DB) error {
	// read schema file
	schema, err := fs.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// execute schema
	if _, err := db.Exec(string(schema)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}
	return nil
}
