package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// CreatePost creates a new post
func CreatePost(db *sql.DB, post *Post) error {
	if err := ValidatePost(post); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	post.ID = uuid.New()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	_, err = tx.Exec(`
	INSERT INTO posts (id, user_id, title, content, created_at, updated_at, category_id)
	VALUES (?, ?, ?, ?, ?, ?, ?)	
	`, post.ID, post.UserID, post.Title, post.Content, post.CreatedAt, post.UpdatedAt, post.Category.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
