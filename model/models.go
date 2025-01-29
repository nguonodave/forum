package models

import (
	"time"

	"github.com/google/uuid"
)

// user struct
type User struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password string
}

// post struct
type Post struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	CategoryID uuid.UUID
	Title      string
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time

	// relationships
	User     *User
	Category *Category
	Comments []*Comment
	Votes    int
}

// Category struct
type Category struct {
	ID        uuid.UUID
	PostId    uuid.UUID
	Name      string
	CreatedAt time.Time
}

// comment struct
type Comment struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PostID    uuid.UUID
	ParentID  *uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      *User
}

// votes struct(likes and dislikes)
type Vote struct {
	UserID    uuid.UUID
	PostID    uuid.UUID
	Value     int
	CreatedAt time.Time
}

// session struct
type Session struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
}
