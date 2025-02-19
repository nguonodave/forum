package controller

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"forum/database"
	"forum/model"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // SQL driver
)

// generateSessionToken generates a unique session token using UUID
func generateSessionToken() string {
	return uuid.New().String()
}

func HandleRegister(username, email, password string) error {
	if err := model.ValidateEmail(email); err != nil {
		return err
	}

	if err := model.ValidatePassword(password); err != nil {
		return err
	}

	if model.IsEmailTaken(email) {
		return errors.New("email is already taken")
	}

	if model.IsUserNameTaken(username) {
		return errors.New("username is already taken")
	}

	hashedPassword, err := model.HashPassword(password)
	if err != nil {
		return errors.New("internal server error")
	}

	userID := uuid.New().String()

	_, DBErr := database.Db.Exec(
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

	fmt.Printf("user %s was created successfully\n", username)
	return nil
}

func HandleLogin(email, username, password string) (string, time.Time, error) {
	var user model.User
	if email == "" && username == "" {
		return "", time.Time{}, errors.New("email and username is missing")
	}

	switch {
	case email != "":

		err := database.Db.QueryRow(
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
		row := database.Db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username)
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

	// Remove any existing sessions for this user before creating a new one
	_, err := database.Db.Exec("DELETE FROM sessions WHERE user_id = ?", user.ID)
	if err != nil {
		log.Println("ERROR:", err)
		return "", time.Time{}, errors.New("internal server error")
	}

	sessionToken := generateSessionToken()

	expiresAt := time.Now().Add(24 * 14 * time.Hour)
	_, err = database.Db.Exec(
		"INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)",
		user.ID, sessionToken, expiresAt,
	)
	if err != nil {
		return "", time.Time{}, errors.New("internal server error")
	}

	return sessionToken, expiresAt, nil
}

// ValidateSession applies for routes that require authentication
func ValidateSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var userID string
		var expiresAt time.Time
		var username string

		row := database.Db.QueryRow("SELECT user_id, expires_at expires_at FROM sessions WHERE token = ? LIMIT 1",
			cookie.Value,
		)

		err = row.Scan(&userID, &expiresAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			fmt.Printf("ERROR: failed to scan session: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if time.Now().After(expiresAt) {
			http.Error(w, "session expired", http.StatusUnauthorized)
			return
		}
		fmt.Println("user id", userID)
		err = database.Db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
		if err != nil {
			log.Println("ERROR: failed to get username")
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "userId", userID)
		ctx = context.WithValue(ctx, "username", username)
		next(w, r.WithContext(ctx))
	}
}

func GetLikesDislikesForPost(db *sql.DB, postId string, Likes *int, Dislikes *int) error {
	query := `
        SELECT 
            COALESCE(SUM(CASE WHEN type = 'like' THEN 1 ELSE 0 END), 0) AS likes,
            COALESCE(SUM(CASE WHEN type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes
        FROM votes
        WHERE post_id = ?
    `

	// Scan directly into the provided pointers
	err := db.QueryRow(query, postId).Scan(Likes, Dislikes)
	if err != nil {
		return err
	}

	return nil
}

// Get likes and dislikes for a certain comment and also return an error
func GetLikesDislikesForComment(db *sql.DB, userId, commentId string) (int, int, error) {
	query := `
        SELECT 
            SUM(CASE WHEN type = 'like' THEN 1 ELSE 0 END) AS likes,
            SUM(CASE WHEN type = 'dislike' THEN 1 ELSE 0 END) AS dislikes
        FROM votes
        WHERE user_id = ? AND  = ?
    `

	var likes, dislikes sql.NullInt64
	row := db.QueryRow(query, userId, commentId)
	err := row.Scan(&likes, &dislikes)
	if err != nil {
		return 0, 0, err
	}

	// Handle NULL values by converting them to 0
	return int(likes.Int64), int(dislikes.Int64), nil
}
