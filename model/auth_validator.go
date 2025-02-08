package model

// Handles database interaction and business logic.
import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"unicode"

	"forum/xerrors"

	"golang.org/x/crypto/bcrypt"
)

var (
	Cost                  int = bcrypt.DefaultCost
	MinimumPasswordLength int = 8
)

// IsValidPassword compares the password and hashedPassword and checks if they match, if not it returns False else True (meaning they match)
func IsValidPassword(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// HashPassword attempts to hash password using Cost value and returns the hashed password and error will be nil if successful
// else hashed password will be an empty string and error will be not nil
func HashPassword(password string) (string, error) {
	fmt.Println("hashA")
	if err := ValidatePassword(password); err != nil {
		fmt.Println("hashB")
		return "", xerrors.ErrPasswordTooShort
	}
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), Cost)
	if err != nil {
		fmt.Println("hashC")
		return "", err
	}
	fmt.Println("hashD")
	return string(passwordBytes), nil
}

func ValidatePassword(password string) error {
	if len(password) < MinimumPasswordLength {
		return xerrors.ErrPasswordTooShort
	}

	if len(password) > 64 {
		return errors.New("password too long")
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
	// Improved regex pattern
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailPattern)
	match := re.MatchString(email)
	fmt.Println("email matches regex pattern", match, email)
	if !match {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// IsEmailTaken queries the database to check if the email provided exists returns true if found else false
func IsEmailTaken(DBase *Database, email string) bool {
	if DBase.Db == nil {
		log.Println("IsEmailTaken() received a nil database connection")
		return false
	}
	var emailExists bool
	println()
	println()
	println()
	fmt.Println("IsEmailTaken() function failure")
	err := DBase.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&emailExists)
	fmt.Println(err)
	if err != nil {
		fmt.Printf("Error checking if email exists: %v\n", err)
		return false
	}
	return emailExists
}
