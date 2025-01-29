package models

// Handles database interaction and business logic.
import (
	"forum/xerrors"

	"golang.org/x/crypto/bcrypt"
)

var (
	Cost int = bcrypt.DefaultCost
)

// IsValidPassword compares the password and hashedPassword and checks if they match, if not it returns False else True (meaning they match)
func IsValidPassword(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// HashPassword attempts to hash password using Cost value and returns the hashed password and error will be nil if successful
// else hashed password will be an empty string and error will be not nil
func HashPassword(password string) (string, error) {
	if len(password) > 72 {
		return "", bcrypt.ErrPasswordTooLong
	}
	if len(password) < 8 {
		return "", xerrors.ErrPasswordTooShort
	}
	password_bytes, err := bcrypt.GenerateFromPassword([]byte(password), Cost)
	if err != nil {
		return "", err
	}
	return string(password_bytes), nil
}
