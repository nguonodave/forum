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

// user struct
type User struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
}

// post struct
type Post struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Title     string
	Content   string
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

// comment struct
type Comment struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	UserID    uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      *User
}

// votes struct(likes and dislikes)
type Vote struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PostID    uuid.UUID
	CommentID uuid.UUID
	Type      string
	CreatedAt time.Time
}

// session struct
type Session struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
}
