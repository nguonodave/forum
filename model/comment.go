package model

import (
	"database/sql"

	"github.com/google/uuid"
)

// GetCommentsByPostID retrieves all comments for a post
func GetCommentsByPostID(db *sql.DB, postID uuid.UUID) ([]*Comment, error) {
	rows, err := db.Query(`
		SELECT 
			c.id, c.user_id, c.content, c.created_at,
			u.id, u.username
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{
			User: &User{},
		}
		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.User.ID,
			&comment.User.Username,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
