package controller

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"forum/model"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // SQL driver
)

// generateSessionToken generates a unique session token using UUID
func generateSessionToken() string {
	return uuid.New().String()
}

func HandleRegister(DBase *model.Database, username, email, password string) error {
	if DBase == nil || DBase.Db == nil {
		return errors.New("database connection is missing")
	}

	if err := model.ValidateEmail(email); err != nil {
		return err
	}

	if err := model.ValidatePassword(password); err != nil {
		return err
	}

	if model.IsEmailTaken(DBase, email) {
		return errors.New("email is already taken")
	}

	if model.IsUserNameTaken(DBase, username) {
		return errors.New("username is already taken")
	}

	hashedPassword, err := model.HashPassword(password)
	if err != nil {
		return errors.New("internal server error")
	}

	userID := uuid.New().String()

	_, DBErr := DBase.Db.Exec(
		"INSERT INTO users (id, email, password, username) VALUES (?, ?, ?, ?);",
		userID,
		email,
		hashedPassword,
		username,
	)
	if DBErr != nil {
		if strings.Contains(DBErr.Error(), "UNIQUE constraint failed") {
			return errors.New("email or username already exists")
		}
		return errors.New("failed to create user")
	}

	fmt.Println("Success: User created!")
	return nil
}

func HandleLogin(DBase *model.Database, email, username, password string) (string, time.Time, error) {

	var user model.User
	if email == "" && username == "" {
		return "", time.Time{}, errors.New("email and username is missing")
	}

	switch {
	case email != "":

		err := DBase.Db.QueryRow(
			"SELECT id, email, password FROM users WHERE email = ?",
			email,
		).Scan(&user.ID, &user.Email, &user.Password)

		if errors.Is(err, sql.ErrNoRows) {
			return "", time.Time{}, errors.New("invalid credentials")
		}

		if err != nil {
			return "", time.Time{}, errors.New("internal server error")
		}
	case username != "":
		row := DBase.Db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username)
		err := row.Scan(&user.ID, &user.Username, &user.Password)
		if errors.Is(err, sql.ErrNoRows) {
			return "", time.Time{}, errors.New("invalid credentials")
		}
		if err != nil {
			return "", time.Time{}, errors.New("internal server error")
		}
	}

	if ok := model.IsValidPassword(password, user.Password); !ok {
		return "", time.Time{}, errors.New("invalid credentials")
	}

	sessionToken := generateSessionToken()

	expiresAt := time.Now().Add(24 * 14 * time.Hour)
	_, err := DBase.Db.Exec(
		"INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)",
		user.ID, sessionToken, expiresAt,
	)
	if err != nil {
		return "", time.Time{}, errors.New("internal server error")
	}

	return sessionToken, expiresAt, nil
}

// ValidateSession applies for routes that require authentication
func ValidateSession(DBase *model.Database, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var userID string
		var expiresAt time.Time

		row := DBase.Db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE token = ? LIMIT 1",
			cookie.Value,
		)

		err = row.Scan(&userID, &expiresAt)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			fmt.Printf("ERROR: failed to scan session: %v\n", err)
			fmt.Println("128", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if time.Now().After(expiresAt) {
			http.Error(w, "session expired", http.StatusUnauthorized)
			return
		}

		// store user id in context for next handlers
		ctx := context.WithValue(r.Context(), "userID", userID)
		// if you want to retrieve user id in the next handlers use the syntax below:
		// userID = r.Context().Value("userID").(string)
		next(w, r.WithContext(ctx))
	}
}
