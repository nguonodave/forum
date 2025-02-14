package model

import (
	"forum/database"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// GetCategories fetches categories from the database
func GetCategories() ([]Category, error) {
	query := `SELECT id, name FROM categories;`

	rows, err := database.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var cat Category
		var idStr string

		if err := rows.Scan(&idStr, &cat.Name); err != nil {
			return nil, err
		}

		// Convert string ID to uuid.UUID
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		cat.ID = id

		categories = append(categories, cat)
	}

	return categories, nil
}
