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

// GetPostByID retrieves a post by its ID including relationships
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

func GetPostsByCategory(db *sql.DB, categoryID uuid.UUID) ([]*Post, error) {
	rows, err := db.Query(`
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at, p.updated_at,
			u.id, u.username,
			COALESCE(v.vote_count, 0) as votes
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN (
			SELECT post_id, SUM(value) as vote_count 
			FROM votes 
			GROUP BY post_id
		) v ON p.id = v.post_id
		WHERE p.category_id = ?
		ORDER BY p.created_at DESC
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{
			User:     &User{},
			Category: &Category{ID: categoryID},
		}
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.User.ID,
			&post.User.Username,
			&post.Votes,
		)
		if err != nil {
			return nil, err
		}

		// Get comments for each post
		comments, err := GetCommentsByPostID(db, post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments

		posts = append(posts, post)
	}

	return posts, nil
}

// UpdatePost updates an existing post
func UpdatePost(db *sql.DB, post *Post) error {
	if err := ValidatePost(post); err != nil {
		return err
	}

	post.UpdatedAt = time.Now()

	result, err := db.Exec(`
	
	UPDATE posts
	SET title = ?, content = ?, updated_at = ?, category_id = ?
	WHERE id = ? AND user_id = ?
	`, post.Title, post.Content, post.UpdatedAt, post.Category.ID, post.ID, post.UserID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return xerrors.ErrInvalidPost
	}

	return nil
}

// DeletePost deletes a post and its associated comments
func DeletePost(db *sql.DB, postID, userID uuid.UUID) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete associated comments first
	_, err = tx.Exec("DELETE FROM comments WHERE post_id = ?", postID)
	if err != nil {
		return err
	}

	// Delete associated votes
	_, err = tx.Exec("DELETE FROM votes WHERE post_id = ?", postID)
	if err != nil {
		return err
	}

	// Delete the post
	result, err := tx.Exec("DELETE FROM posts WHERE id = ? AND user_id = ?", postID, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return xerrors.ErrInvalidPost
	}

	return tx.Commit()
}

// ValidatePost checks if any required section of a post is missing
// It returns an error if any section is missing
func ValidatePost(post *Post) error {
	if post.Title == "" {
		return xerrors.ErrEmptyTitle
	}

	if post.Content == "" {
		return xerrors.ErrEmptyContent
	}

	if post.UserID == uuid.Nil {
		return xerrors.ErrInvalidUser
	}

	return nil
}
