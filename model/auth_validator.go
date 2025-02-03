package model

// Handles database interaction and business logic.
import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"unicode"

	"forum/xerrors"

	"golang.org/x/crypto/bcrypt"
)

var (
	Cost                  int = bcrypt.DefaultCost
	MinimumPasswordLength int = 8
)

// VerifyPassword compares the password and hashedPassword and checks if they match, if not it returns False else True (meaning they match)
func VerifyPassword(password, hashedPassword string) bool {
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

// ValidateEmail checks if email provided has a valid email syntax
func ValidateEmail(email string) error {
	emailPattern := `^(\w)(\w{1,}|\d{1,})+@\w+.\w{2,4}$`
	if !regexp.MustCompile(emailPattern).MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// IsEmailTaken queries the database to check if the email provided exists returns true if found else false
func IsEmailTaken(email string, db *sql.DB) bool {
	var emailExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&emailExists)
	if err != nil {
		fmt.Printf("Error checking if email exists: %v\n", err)
		return false
	}
	return emailExists
}
