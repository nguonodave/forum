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

func HandleRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		Email = r.FormValue("email")
		Password = r.FormValue("password")
		Username = r.FormValue("username")

		// Validate input
		if err := model.ValidateEmail(Email); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := model.ValidatePassword(Password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if model.IsEmailTaken(db, Email) {
			http.Error(w, "Email already registered", http.StatusConflict)
			return
		}

		// Hash password
		hashedPassword, err := model.HashPassword(Password)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Insert user into database
		_, err = db.Exec(
			"INSERT INTO users (email, password, username) VALUES (?, ?, ?)",
			Email,
			hashedPassword,
			Username,
		)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func VerifyLogin(email, password, username string) (string, string, error) {
	// Retrieve user from database
	var user model.User
	err := db.QueryRow(
		"SELECT id, email, password FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &user.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return "", "", errors.New("Invalid credentials")
	}

	if err != nil {
		return "", "", errors.New("Internal server error")
	}

	// Verify password
	if ok := model.VerifyPassword(user.Password, password); ok != true {
		return "", "", errors.New("Invalid credentials")
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
