package models

import (
	"time"

	"github.com/google/uuid"
)

// user struct
type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username" validate: "required,min=3"`
	Email     string    `json:"email" validate: "required, email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// post struct
type Post struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json: "user_id"`
	CategoryID uuid.UUID `json:"category_id"`
	Title      string    `json:"title" validate:"required"`
	Content    string    `json:"content" validate: "required"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json: "deleted_at"`

	// relationships
	User     *User      `json:"user, omitempty"`
	Category *Category  `json: "category, omitempty"`
	Comments []*Comment `json: "comments, omitempty"`
	Votes    int        `json:"votes"`
}

// Category struct
type Category struct {
	ID        uuid.UUID `json:"id"`
	PostId    uuid.UUID `json: "post_id"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

// comment struct
type Comment struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	PostID    uuid.UUID  `json:"post_id"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	Content   string     `json:"content" validate:"required,min=2"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Relationships
	User *User `json:"user,omitempty"`
	// Replies []*Comment `json:"replies,omitempty"`
}

// votes struct(likes and dislikes)
type Vote struct {
	UserID    uuid.UUID `json:"-"`
	PostID    uuid.UUID `json:"post_id"`
	Value     int       `json:"value" validate:"oneof=-1 1"`
	CreatedAt time.Time `json:"created_at"`
}


//session struct
type Session struct {
	Token string `json:"-"`
	UserID uuid.UUID `json:"-"`
	ExpiresAt time.Time `json: "expires_at"`
}