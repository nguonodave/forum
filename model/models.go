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
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}
