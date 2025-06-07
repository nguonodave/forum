package controllers

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetUserId using username or email
func GetUserId(identifier string) (int, error) {
	query := `SELECT user_id FROM users WHERE name = $1 OR email = $1 LIMIT 1`
	var userID int
	err := DB.QueryRow(query, identifier).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user not found")
		}
		return 0, err
	}
	return userID, nil
}
