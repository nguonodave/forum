package model

// Handles database interaction and business logic.
import (
	"errors"
	"forum/xerrors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)


var (
	Cost                  int = bcrypt.DefaultCost
	MinimumPasswordLength int = 8
)

// VerifyPassword compares the password and hashedPassword and checks if they match, if not it returns False else True (meaning they match)
func IsValidPassword(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// HashPassword attempts to hash password using Cost value and returns the hashed password and error will be nil if successful
// else hashed password will be an empty string and error will be not nil
func HashPassword(password string) (string, error) {
	if err := ValidatePassword(password); err != nil {
		return "", xerrors.ErrPasswordTooShort
	}
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), Cost)
	if err != nil {
		return "", err
	}
	return string(passwordBytes), nil
}

func ValidatePassword(password string) error {
	if len(password) < MinimumPasswordLength {
		return xerrors.ErrPasswordTooShort
	}

	var hasUpper, hasLower, hasNumber, hasPunct bool

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c):
			hasPunct = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	if !hasPunct {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
