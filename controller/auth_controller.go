package controller

import (
	"context"
	"database/sql"
	"errors"
	"forum/xerrors"
	"log"
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
	DB       *sql.DB
)

// generateSessionToken generates a unique session token using UUID
func generateSessionToken() string {
	return uuid.New().String()
}

func HandleRegister(username, email, password) error {
	if err := model.ValidateEmail(Email); err != nil {
		return err
	}

	if err := model.ValidatePassword(Password); err != nil {
		return err
	}

	if model.IsEmailTaken(DB, Email) {
		return errors.New("email is already taken")
	}

	hashedPassword, err := model.HashPassword(Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	_, err = DB.Exec(
		"INSERT INTO users (email, password, username) VALUES (?, ?, ?)",
		Email,
		hashedPassword,
		Username,
	)

	if err != nil {
		return err
	}
	return nil
}

// HandleLogin takes only the email and password and creates and checks if user exist in database
// it returns:
// 1. sessionToken else empty string
// 2. time the session token expires else empty time struct
// 3. error == nil if success
func HandleLogin(email, password string) (string, time.Time, error) {
	var user model.User
	err := DB.QueryRow(
		"SELECT id, email, password FROM users WHERE email = ?",
		Email,
	).Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", time.Time{}, xerrors.ErrNoSuchUser
		}
		return "", time.Time{}, err
	}

	if ok := model.VerifyPassword(user.Password, Password); ok != true {
		return "", time.Time{}, xerrors.ErrInvalidCredentials
	}

	sessionToken := generateSessionToken()

	// expire after a day
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = DB.Exec(
		"INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)",
		sessionToken,
		user.ID,
		expiresAt,
	)
	return sessionToken, expiresAt, nil
}

// ValidateSession is a middleware that checks if a users session exist in database and has not expired
// to be used in routes that need authentication
func ValidateSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the session cookie
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var userID int
		var expiresAt time.Time

		err = DB.QueryRow(
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
		_, err = DB.Exec("UPDATE sessions SET expires_at = ? WHERE token = ?", time.Now().Add(24*time.Hour), cookie.Value)
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
