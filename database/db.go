package database

import (
	"database/sql"
	"embed"
	"fmt"
	"forum/model"

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

func InitializeDB() (*model.Database, error) {
	// Open database connection
	dsn := fmt.Sprintf("file:%s?_foreign_keys=on&_journal_mode=WAL", dbName)
	sqlDB, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Wrap in model.Database
	db := &model.Database{Db: sqlDB}

	// Apply migrations
	if err := applyMigration(db); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	return db, nil // ✅ Return the initialized database
}

// Apply migrations to the database.
func applyMigration(db *model.Database) error {
	// read schema file
	schema, err := fs.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// execute schema
	if _, err := db.Db.Exec(string(schema)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}
	return nil
}
