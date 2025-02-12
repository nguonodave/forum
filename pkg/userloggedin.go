package pkg

import (
	"database/sql"
	"net/http"
	"time"
)

func UserLoggedIn(r *http.Request, db *sql.DB) bool {
	// get cookie from request
	cookie, cookieErr := r.Cookie("session")
	if cookieErr != nil {
		return false
	}

	// check if cookie details are in sessions table
	var userId string
	var expiryDate time.Time
	sessionQueryErr := db.QueryRow(`SELECT user_id, expires_at FROM sessions WHERE id = ?`, cookie.Value).Scan(&userId, &expiryDate)

	// if no err, meaning session is available, return true
	return sessionQueryErr == nil && time.Now().Before(expiryDate)
}
