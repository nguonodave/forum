package middlewares

import (
	"database/sql"
	"net/http"

	"forum/pkg"
)

func RedirectIfLoggedIn(db *sql.DB, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// use UserLoggedIn from pkg to check if user is logged in
		// if so redirect to home
		if pkg.UserLoggedIn(r, db) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// if not logged in, proceed to the next handler, Login
		handler(w, r)
	}
}
