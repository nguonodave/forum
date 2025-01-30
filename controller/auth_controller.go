package controller

import (
	"database/sql"
	"net/http"
	"time"

	"forum/model"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // SQL driver
)

var (
	Email    string
	Password string
	Username string
)

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

func HandleLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		Email = r.FormValue("email")
		Password = r.FormValue("password")
		Username = r.FormValue("username")

		// Retrieve user from database
		var user model.User
		err := db.QueryRow(
			"SELECT id, email, password FROM users WHERE email = ?",
			Email,
		).Scan(&user.ID, &user.Email, &user.Password)

		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Verify password
		if err := model.VerifyPassword(user.Password, Password); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Generate session token
		sessionToken := generateSessionToken()

		// Store session in the database
		expiresAt := time.Now().Add(24 * time.Hour)
		_, err = db.Exec(
			"INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)",
			sessionToken,
			user.ID,
			expiresAt,
		)
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		// Create session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    sessionToken,
			Expires:  expiresAt,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})

		w.WriteHeader(http.StatusOK)
	}
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

		if err == sql.ErrNoRows {
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
