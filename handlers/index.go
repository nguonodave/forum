package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

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
		templateFile := "base.html"
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
		temp, err := template.ParseFiles(filepath.Join(templatesDir, templateFile))
		if err != nil {
			TemplateError("error parsing template", err)
			return
		}

		err = temp.Execute(w, content)
		if err != nil {
			TemplateError("error executing template", err)
			return
		}
	}
}
