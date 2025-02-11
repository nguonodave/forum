package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Database struct
type Database struct {
	Db *sql.DB
}

// User struct
type User struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
}

// Post struct
type Post struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Title     string
	Content   string
	ImageURL  string
	CreatedAt time.Time
	UpdatedAt time.Time

	// relationships
	User     *User
	Category *Category
	Comments []*Comment
	Votes    int
}

// Category struct
type Category struct {
	ID   *uuid.UUID
	Name *string
}

// Comment struct
type Comment struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	UserID    uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      *User
}

// Vote struct(likes and dislikes)
type Vote struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PostID    uuid.UUID
	CommentID uuid.UUID
	Type      string
	CreatedAt time.Time
}

// Session struct
type Session struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
}
