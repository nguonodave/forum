package models

import (
	"time"

	"github.com/google/uuid"
)

// user struct
type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username" validate: "required, min=3"`
	Email     string    `json:"email" validate: "required, email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
