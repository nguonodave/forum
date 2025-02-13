package database

import (
	"database/sql"
	"embed"
	"fmt"

	"forum/model"

	"github.com/google/uuid"
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

	return db, nil
}

// Apply migrations to the database.
func applyMigration(db *model.Database) error {
	// Read schema file
	schema, err := fs.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute schema
	if _, err := db.Db.Exec(string(schema)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	// Insert predefined categories
	if err := insertDefaultCategories(db); err != nil {
		return fmt.Errorf("failed to insert default categories: %w", err)
	}

	return nil
}

// Insert default categories into the database.
func insertDefaultCategories(db *model.Database) error {
	categories := []string{"General", "Technology", "Sports", "Entertainment", "Health"}

	for _, name := range categories {
		id := uuid.New() // Generate a new UUID

		_, err := db.Db.Exec(
			"INSERT INTO categories (id, name) VALUES (?, ?) ON CONFLICT(name) DO NOTHING;",
			id.String(), name,
		)
		if err != nil {
			return fmt.Errorf("failed to insert category %s: %w", name, err)
		}
	}

	return nil
}
