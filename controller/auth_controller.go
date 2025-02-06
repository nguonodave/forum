package controller

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"forum/model"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // SQL driver
)

// Expose the db variable for testing
var db *sql.DB

// SetDB sets the global database connection
func SetDB(database *sql.DB) {
	db = database
}

// generateSessionToken generates a unique session token using UUID
func generateSessionToken() string {
	return uuid.New().String()
}

// HandleRegister registers a new user
func HandleRegister(username, email, password string) error {
	// Validate input
	println(333)
	if err := model.ValidateEmail(email); err != nil {
		println(3334)
		return err
	}

	println(3335)
	if err := model.ValidatePassword(password); err != nil {
		println(3336)
		return err
	}

	println(3337)
	//if model.IsEmailTaken(db, email) {
	//	println(3338)
	//	return errors.New("email is already taken")
	//}

	// Hash password
	hashedPassword, err := model.HashPassword(password)
	fmt.Println(hashedPassword)
	println(33322)
	if err != nil {
		println(3339)
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
		println(33310)
		return errors.New("failed to create user")
	}

	println(33311)
	fmt.Println("success in creating user")
	return nil
}

func HandleLogin(email, password string) (string, time.Time, error) {
	// Retrieve user from database
	var user model.User
	err := db.QueryRow(
		"SELECT id, email, password FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &user.Password)

	if errors.Is(err, sql.ErrNoRows) {
		return "", time.Time{}, errors.New("invalid credentials")
	}

	if err != nil {
		return "", time.Time{}, errors.New("internal server error")
	}

	// Verify password
	if ok := model.IsValidPassword(password, user.Password); !ok {
		return "", time.Time{}, errors.New("invalid credentials")
	}

	// Generate session token
	sessionToken := generateSessionToken()

	// Store session in the database
	expiresAt := time.Now().Add(24 * 14 * time.Hour)
	_, err = db.Exec(
		"INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)",
		user.ID, sessionToken, expiresAt,
	)
	if err != nil {
		return "", time.Time{}, errors.New("internal server error")
	}

	return sessionToken, expiresAt, nil
}

// ValidateSession applies for routes that require authentication
func ValidateSession(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the session cookie
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var userID int
		var expiresAt time.Time

		err = db.QueryRow(
			"SELECT user_id, expires_at FROM sessions WHERE token = ? LIMIT 1",
			cookie.Value,
		).Scan(&userID, &expiresAt)

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err != nil {
			log.Printf("ERROR: database while validating session %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if time.Now().After(expiresAt) {
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}

		// OPTIONAL FEATURE: refreshing token expiration for each request
		// this is to make sure if the site is idle, we log out user to save resources
		_, err = db.Exec("UPDATE sessions SET expires_at = ? WHERE token = ?", time.Now().Add(24*time.Hour), cookie.Value)
		if err != nil {
			log.Printf("ERROR: database while refershing session %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// store user id in context for next handlers
		ctx := context.WithValue(r.Context(), "userID", userID)
		// if you want to retrieve user id in the next handlers use the syntax below:
		// userID = r.Context().Value("userID").(int)
		next(w, r.WithContext(ctx))
	}
}
