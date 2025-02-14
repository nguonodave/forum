package database

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// Database connection and initialization.
//
//go:embed schema.sql
var fs embed.FS
var Db *sql.DB

const (
	dbName     = "forum.db"
	schemaFile = "schema.sql"
)

func InitializeDB() (error) {
	// Open database connection
	dsn := fmt.Sprintf("file:%s?_foreign_keys=on&_journal_mode=WAL", dbName)
	var openDbErr error
	Db, openDbErr = sql.Open("sqlite3", dsn)
	if openDbErr != nil {
		return fmt.Errorf("failed to open database: %w", openDbErr)
	}

	// Verify connection
	if err := Db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Apply migrations
	if err := applyMigration(); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Configure connection pool
	Db.SetMaxOpenConns(1)
	Db.SetMaxIdleConns(1)

	return nil
}

// Apply migrations to the database.
func applyMigration() error {
	// Read schema file
	schema, err := fs.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute schema
	if _, err := Db.Exec(string(schema)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	// Insert predefined categories
	if err := insertDefaultCategories(); err != nil {
		return fmt.Errorf("failed to insert default categories: %w", err)
	}

	return nil
}

// Insert default categories into the database.
func insertDefaultCategories() error {
	categories := []string{"General", "Technology", "Sports", "Entertainment", "Health"}

	for _, name := range categories {
		id := uuid.New() // Generate a new UUID

		_, err := Db.Exec(
			"INSERT INTO categories (id, name) VALUES (?, ?) ON CONFLICT(name) DO NOTHING;",
			id.String(), name,
		)
		if err != nil {
			return fmt.Errorf("failed to insert category %s: %w", name, err)
		}
	}

	return nil
}
