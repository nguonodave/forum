package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"forum/model"

	_ "github.com/mattn/go-sqlite3" // SQL driver
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username,omitempty"`
}

func HandleRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate input
		if err := model.ValidateEmail(creds.Email); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := model.ValidatePassword(creds.Password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if model.IsEmailTaken(db, creds.Email) {
			http.Error(w, "Email already registered", http.StatusConflict)
			return
		}

		// Hash password
		hashedPassword, err := model.HashPassword(creds.Password)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Insert user into database
		_, err = db.Exec(
			"INSERT INTO users (email, password, username) VALUES (?, ?, ?)",
			creds.Email,
			hashedPassword,
			creds.Username,
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
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Retrieve user from database
		var user model.User
		err := db.QueryRow(
			"SELECT id, email, password FROM users WHERE email = ?",
			creds.Email,
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
		if err := model.VerifyPassword(user.Password, creds.Password); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Create session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    generateSessionToken(), // pending implementation
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})

		w.WriteHeader(http.StatusOK)
	}
}
