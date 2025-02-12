package middlewares

import (
	"database/sql"
	"net/http"

	"forum/model"
	"forum/pkg"
)

func RedirectIfLoggedIn(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// use UserLoggedIn from pkg to check if user is logged in
		// if so redirect to home
		var db *sql.DB
		database := model.Database{Db: db}
		if pkg.UserLoggedIn(r, database.Db) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// if not logged in, proceed to the next handler, Login
		handler(w, r)
	}
}
