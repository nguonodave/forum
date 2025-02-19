package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"forum/pkg"
)

func RedirectIfLoggedIn(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// use UserLoggedIn from pkg to check if user is logged in
		// if so redirect to home
		isLoggedIn, username := pkg.UserLoggedIn(r)
		fmt.Println("???", isLoggedIn, username)
		if isLoggedIn {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), "username", username)
		r = r.WithContext(ctx)
		// if not logged in, proceed to the next handler, Login
		handler(w, r)
	}
}
