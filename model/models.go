package model

import (
	"time"

	"github.com/google/uuid"
)

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

	// Relationships
	User       *User
	Categories []*Category
	Comments   []*Comment
	Votes      []*Vote // Changed from int to slice to store actual votes
}

// Category struct
type Category struct {
	ID   uuid.UUID
	Name string
}

// Comment struct
type Comment struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	UserID    uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	User  *User
	Votes []*Vote // Enables counting likes/dislikes dynamically
}

// Vote struct
type Vote struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PostID    *uuid.UUID // Nullable, as it may belong to a post or comment
	CommentID *uuid.UUID // Nullable, as it may belong to a comment
	Type      string     // "like" or "dislike"
	CreatedAt time.Time
}

// Session struct
type Session struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
}
