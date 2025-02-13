// utils/session.go
package helperfunc

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"forum/model"

	"github.com/google/uuid"
)

// Session cookie name constant
const SessionCookieName = "session"

// GetSessionToken retrieves the session token from the request cookie
func GetSessionToken(r *http.Request) (string, error) {
	// Get session cookie
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// GetUserIDFromSession retrieves the user ID associated with the current session
func GetUserIDFromSession(r *http.Request) (uuid.UUID, error) {
	// Get session token
	token, err := GetSessionToken(r)
	if err != nil {
		return uuid.Nil, err
	}

	// Get database connection from context
	db := r.Context().Value("db").(*model.Database)
	if db == nil {
		return uuid.Nil, errors.New("database connection not found in context")
	}

	// Query the database for session info
	var userID uuid.UUID
	var expiresAt time.Time

	err = db.Db.QueryRow(`
		SELECT user_id, expires_at 
		FROM sessions 
		WHERE token = ?
	`, token).Scan(&userID, &expiresAt)

	if err == sql.ErrNoRows {
		return uuid.Nil, errors.New("invalid session")
	}
	if err != nil {
		return uuid.Nil, err
	}

	// Check if session has expired
	if time.Now().After(expiresAt) {
		// Delete expired session
		_, _ = db.Db.Exec("DELETE FROM sessions WHERE token = ?", token)
		return uuid.Nil, errors.New("session expired")
	}

	return userID, nil
}
