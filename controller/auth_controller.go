package controller

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"forum/model"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // SQL driver
)

var db *sql.DB

// generateSessionToken generates a unique session token using UUID
func generateSessionToken() string {
	return uuid.New().String()
}

func RegisterUser(email, password, username string) error {

	// Validate input
	if err := model.ValidateEmail(email); err != nil {

		return err
	}

	if err := model.ValidatePassword(password); err != nil {

		return err
	}

	if model.IsEmailTaken(db, email) {

		return errors.New("email is already taken")
	}

	// Hash password
	hashedPassword, err := model.HashPassword(password)
	if err != nil {
		return errors.New("internal server error")
	}

	// Insert user into database
	_, err = db.Exec(
		"INSERT INTO users (email, password, username) VALUES (?, ?, ?)",
		email,
		hashedPassword,
		username,
	)
	if err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func VerifyLogin(email, password, username string) (string, string, error) {
	// Retrieve user from database
	var user model.User
	err := db.QueryRow(
		"SELECT id, email, password FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &user.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return "", "", errors.New("invalid credentials")
	}

	if err != nil {
		return "", "", errors.New("internal server error")
	}

	// Verify password
	if ok := model.IsValidPassword(user.Password, password); !ok {
		return "", "", errors.New("invalid credentials")
	}

	// Generate session token
	sessionToken := generateSessionToken()

	// Store session in the database
	expiresAt := time.Now().Add(24 * time.Hour)

	return sessionToken, expiresAt.String(), nil
}

// ValidateSession applies for routes that require authentication.
func ValidateSession(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the session cookie
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the session exists in the database
		var userID int
		var expiresAt time.Time
		err = db.QueryRow(
			"SELECT user_id, expires_at FROM sessions WHERE token = ?",
			cookie.Value,
		).Scan(&userID, &expiresAt)

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Check if the session has expired
		if time.Now().After(expiresAt) {
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}

		// Call the next handler
		next(w, r)
	}
}
