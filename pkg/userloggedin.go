package pkg

import (
	"fmt"
	"forum/database"
	"log"
	"net/http"
	"time"
)

func UserLoggedIn(r *http.Request) (bool, string) {
	// get cookie from request
	cookie, cookieErr := r.Cookie("session")
	if cookieErr != nil {
		return false, ""
	}
	var username string

	// check if cookie details are in sessions table
	var userId string
	var expiryDate time.Time
	sessionQueryErr := database.Db.QueryRow(`SELECT user_id, expires_at FROM sessions WHERE token = ?`, cookie.Value).Scan(&userId, &expiryDate)

	fmt.Println(userId, expiryDate)
	err := database.Db.QueryRow("SELECT username from users WHERE id = ?", userId).Scan(&username)
	if err != nil {
		log.Println(err)
		return false, ""
	}
	log.Println(username, "logged in with user id ==>", userId)
	// if no err, meaning session is available, return true
	return sessionQueryErr == nil && time.Now().Before(expiryDate), username
}
