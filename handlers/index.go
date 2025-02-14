package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"forum/model"
	"forum/pkg"
)

// Index handler designed for the application's index page
func Index(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// fetch all categories

		// if method is POST (filters)
			// form value
				// GetPost(db, value)
					// if value not empty, SELECT from db WHERE value
					// else SELECT (*)
			// populate posts
		// else if GET
			// all posts
			// populate posts

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		type Content struct {
			Message      string
			Data         string
			UserLoggedIn bool
			Posts        []model.Post
			User         model.User
		}

		content := Content{
			Message:      "Some message to pass to template",
			Data:         "Some data to pass to template",
			UserLoggedIn: pkg.UserLoggedIn(r, db),
		}

		TemplateError := func(message string, err error) {
			http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
			log.Printf("%s: %v", message, err)
		}

		execTemplateErr := Templates.ExecuteTemplate(w, "base.html", content)
		if execTemplateErr != nil {
			TemplateError("error executing template", execTemplateErr)
			return
		}
	}
}
