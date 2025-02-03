package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"forum/xerrors"
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

func GetPostByID(db *sql.DB, postID uuid.UUID) (*Post, error) {
	post := &Post{}

	err := db.QueryRow(`
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at, p.updated_at,
			u.id, u.username,
			c.id, c.name,
			COALESCE(v.vote_count, 0) as votes
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN (
			SELECT post_id, SUM(value) as vote_count 
			FROM votes 
			GROUP BY post_id
		) v ON p.id = v.post_id
		WHERE p.id = ?
	`, postID).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.User.ID,
		&post.User.Username,
		&post.Category.ID,
		&post.Category.Name,
		&post.Votes,
	)

	if err == sql.ErrNoRows {
		return nil, xerrors.ErrInvalidPost
	}

	if err != nil {
		return nil, err
	}

	comments, err := GetCommentsByPostID(db, post.ID)
	if err != nil {
		return nil, err
	}
	post.Comments = comments

	return post, nil
}
