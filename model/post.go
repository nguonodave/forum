package model

// import (
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"math/rand"
// 	"time"

// 	"forum/xerrors"

// 	"github.com/google/uuid"
// )

// // CreatePost creates a new post with categories
// func CreatePost(db *Database, post *Post) error {
// 	if err := ValidatePost(post); err != nil {
// 		return err
// 	}

// 	tx, err := db.Db.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	defer tx.Rollback()

// 	post.ID = uuid.New()
// 	post.CreatedAt = time.Now()
// 	post.UpdatedAt = time.Now()

// 	// Insert the post
// 	_, err = tx.Exec(`
// 	INSERT INTO posts (id, user_id, title, content, image_url, created_at, updated_at)
// 	VALUES (?, ?, ?, ?, ?, ?, ?)	
// 	`, post.ID, post.UserID, post.Title, post.Content, post.ImageURL, post.CreatedAt, post.UpdatedAt)
// 	if err != nil {
// 		return err
// 	}

// 	// Insert post categories
// 	for _, category := range post.Categories {
// 		_, err = tx.Exec(`
// 		INSERT INTO post_categories (post_id, category_id)
// 		VALUES (?, ?)
// 		`, post.ID, category.ID)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return tx.Commit()
// }

// // GetPostByID retrieves a post by its ID including relationships
// func GetPostByID(db *Database, postID uuid.UUID) (*Post, error) {
// 	post := &Post{
// 		User:       &User{},
// 		Categories: make([]*Category, 0),
// 		Comments:   make([]*Comment, 0),
// 		Votes:      make([]*Vote, 0),
// 	}

// 	// Get post details and user info
// 	err := db.Db.QueryRow(`
// 		SELECT 
// 			p.id, p.user_id, p.title, p.content, p.image_url, p.created_at, p.updated_at,
// 			u.id, u.username, u.email
// 		FROM posts p
// 		LEFT JOIN users u ON p.user_id = u.id
// 		WHERE p.id = ?
// 	`, postID).Scan(
// 		&post.ID,
// 		&post.UserID,
// 		&post.Title,
// 		&post.Content,
// 		&post.ImageURL,
// 		&post.CreatedAt,
// 		&post.UpdatedAt,
// 		&post.User.ID,
// 		&post.User.Username,
// 		&post.User.Email,
// 	)

// 	if errors.Is(err, sql.ErrNoRows) {
// 		return nil, xerrors.ErrInvalidPost
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Get categories
// 	rows, err := db.Db.Query(`
// 		SELECT c.id, c.name
// 		FROM categories c
// 		JOIN post_categories pc ON c.id = pc.category_id
// 		WHERE pc.post_id = ?
// 	`, postID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		category := &Category{}
// 		if err := rows.Scan(&category.ID, &category.Name); err != nil {
// 			return nil, err
// 		}
// 		post.Categories = append(post.Categories, category)
// 	}

// 	// Get votes
// 	rows, err = db.Db.Query(`
// 		SELECT id, user_id, type, created_at
// 		FROM votes
// 		WHERE post_id = ?
// 	`, postID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		vote := &Vote{PostID: &postID}
// 		if err := rows.Scan(&vote.ID, &vote.UserID, &vote.Type, &vote.CreatedAt); err != nil {
// 			return nil, err
// 		}
// 		post.Votes = append(post.Votes, vote)
// 	}

// 	// Get comments with their votes
// 	comments, err := GetCommentsByPostID(db, post.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	post.Comments = comments

// 	return post, nil
// }

// func GetPostsByCategory(db *Database, categoryID uuid.UUID) ([]*Post, error) {
// 	rows, err := db.Db.Query(`
// 		SELECT DISTINCT
// 			p.id, p.user_id, p.title, p.content, p.image_url, p.created_at, p.updated_at,
// 			u.id, u.username, u.email
// 		FROM posts p
// 		JOIN post_categories pc ON p.id = pc.post_id
// 		LEFT JOIN users u ON p.user_id = u.id
// 		WHERE pc.category_id = ?
// 		ORDER BY p.created_at DESC
// 	`, categoryID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var posts []*Post
// 	for rows.Next() {
// 		post := &Post{
// 			User:       &User{},
// 			Categories: make([]*Category, 0),
// 			Comments:   make([]*Comment, 0),
// 			Votes:      make([]*Vote, 0),
// 		}

// 		err := rows.Scan(
// 			&post.ID,
// 			&post.UserID,
// 			&post.Title,
// 			&post.Content,
// 			&post.ImageURL,
// 			&post.CreatedAt,
// 			&post.UpdatedAt,
// 			&post.User.ID,
// 			&post.User.Username,
// 			&post.User.Email,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Get categories for each post
// 		categories, err := GetPostCategories(db, post.ID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		post.Categories = categories

// 		// Get votes for each post
// 		votes, err := GetPostVotes(db, post.ID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		post.Votes = votes

// 		// Get comments for each post
// 		comments, err := GetCommentsByPostID(db, post.ID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		post.Comments = comments

// 		posts = append(posts, post)
// 	}

// 	return posts, nil
// }

// // UpdatePost updates an existing post
// func UpdatePost(db *Database, post *Post) error {
// 	if err := ValidatePost(post); err != nil {
// 		return err
// 	}

// 	tx, err := db.Db.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	post.UpdatedAt = time.Now()

// 	// Update post
// 	result, err := tx.Exec(`
// 		UPDATE posts
// 		SET title = ?, content = ?, image_url = ?, updated_at = ?
// 		WHERE id = ? AND user_id = ?
// 	`, post.Title, post.Content, post.ImageURL, post.UpdatedAt, post.ID, post.UserID)
// 	if err != nil {
// 		return err
// 	}

// 	rows, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	if rows == 0 {
// 		return xerrors.ErrInvalidPost
// 	}

// 	// Update categories
// 	_, err = tx.Exec("DELETE FROM post_categories WHERE post_id = ?", post.ID)
// 	if err != nil {
// 		return err
// 	}

// 	for _, category := range post.Categories {
// 		_, err = tx.Exec(`
// 			INSERT INTO post_categories (post_id, category_id)
// 			VALUES (?, ?)
// 		`, post.ID, category.ID)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return tx.Commit()
// }

// // DeletePost deletes a post and its associated data
// func DeletePost(db *Database, postID, userID uuid.UUID) error {
// 	tx, err := db.Db.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	// Delete associated votes first (for both post and its comments)
// 	_, err = tx.Exec(`
// 		DELETE FROM votes 
// 		WHERE post_id = ? OR comment_id IN (
// 			SELECT id FROM comments WHERE post_id = ?
// 		)
// 	`, postID, postID)
// 	if err != nil {
// 		return err
// 	}

// 	// Delete associated comments
// 	_, err = tx.Exec("DELETE FROM comments WHERE post_id = ?", postID)
// 	if err != nil {
// 		return err
// 	}

// 	// Delete post categories
// 	_, err = tx.Exec("DELETE FROM post_categories WHERE post_id = ?", postID)
// 	if err != nil {
// 		return err
// 	}

// 	// Delete the post
// 	result, err := tx.Exec("DELETE FROM posts WHERE id = ? AND user_id = ?", postID, userID)
// 	if err != nil {
// 		return err
// 	}

// 	rows, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	if rows == 0 {
// 		return xerrors.ErrInvalidPost
// 	}

// 	return tx.Commit()
// }

// // Helper functions

// func GetPostCategories(db *Database, postID uuid.UUID) ([]*Category, error) {
// 	rows, err := db.Db.Query(`
// 		SELECT c.id, c.name
// 		FROM categories c
// 		JOIN post_categories pc ON c.id = pc.category_id
// 		WHERE pc.post_id = ?
// 	`, postID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var categories []*Category
// 	for rows.Next() {
// 		category := &Category{}
// 		if err := rows.Scan(&category.ID, &category.Name); err != nil {
// 			return nil, err
// 		}
// 		categories = append(categories, category)
// 	}
// 	return categories, nil
// }

// func GetPostVotes(db *Database, postID uuid.UUID) ([]*Vote, error) {
// 	rows, err := db.Db.Query(`
// 		SELECT id, user_id, type, created_at
// 		FROM votes
// 		WHERE post_id = ?
// 	`, postID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var votes []*Vote
// 	for rows.Next() {
// 		vote := &Vote{PostID: &postID}
// 		if err := rows.Scan(&vote.ID, &vote.UserID, &vote.Type, &vote.CreatedAt); err != nil {
// 			return nil, err
// 		}
// 		votes = append(votes, vote)
// 	}
// 	return votes, nil
// }

// // ValidatePost checks if any required section of a post is missing
// func ValidatePost(post *Post) error {
// 	if post.Title == "" {
// 		return xerrors.ErrEmptyTitle
// 	}

// 	if post.Content == "" {
// 		return xerrors.ErrEmptyContent
// 	}

// 	if post.UserID == uuid.Nil {
// 		return xerrors.ErrInvalidUser
// 	}

// 	if len(post.Categories) == 0 {
// 		return xerrors.ErrNoCategory
// 	}

// 	return nil
// }

// // RandomUsername generates a random username using adjectives and nouns and numbers
// func RandomUsername() string {
// 	adjectives := []string{"Brave", "Clever", "Swift", "Mighty", "Silent", "Fierce", "Gentle", "Loyal", "Curious", "Wise", "Mashed", "Tall"}
// 	nouns := []string{"Tiger", "Hawk", "Wolf", "Panther", "Eagle", "Fox", "Lion", "Shark", "Bear", "Otter", "Potato", "Ghost", "Egg"}

// 	source := rand.NewSource(time.Now().UnixNano())
// 	rnd := rand.New(source)

// 	adj := adjectives[rnd.Intn(len(adjectives))]
// 	noun := nouns[rnd.Intn(len(nouns))]
// 	number := rnd.Intn(9000) + 1000

// 	return fmt.Sprintf("%s%s%d", adj, noun, number)
// }
