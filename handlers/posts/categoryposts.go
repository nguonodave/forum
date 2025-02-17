package posts

import (
	"forum/handlers"
	"net/http"
)

func CategoryPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handlers.ErrorPage(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	// queriedCategoryId := r.URL.Query().Get("category")
}
