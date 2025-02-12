package pkg

import (
	"database/sql"
	"net/http"
)

func UserLoggedIn(r *http.Request, db *sql.DB) bool {
	// get cookie from request

	// check if cookie details are in sessions table

	// return true is so
	return true
}
