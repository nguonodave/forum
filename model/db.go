package models

import (
	"database/sql"
	"embed"
	"fmt"
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
}
