package handlers

import (
	"forum/controller"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// templatesDir refers to the filepath of the directory containing the application's templates
var templatesDir = "view"

// Login handles both GET and POST methods, if method is GET it renders the page
// if method is POST it gets the values from the form and internally checks if details exist in the database
func Login(w http.ResponseWriter, r *http.Request) {
	templateFile := "auth/login.html"

	if r.Method == "GET" {
		TemplateError := func(message string, err error) {
			http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
			log.Printf("%s: %v", message, err)
		}
		temp, err := template.ParseFiles(filepath.Join(templatesDir, templateFile))
		if err != nil {
			TemplateError("error parsing template", err)
			return
		}
		err = temp.Execute(w, struct{}{})
		if err != nil {
			TemplateError("error executing template", err)
			return
		}
	} else if r.Method == "POST" {
		// parse form and populate r.Form
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		email := r.Form.Get("email")
		password := r.FormValue("password")
		username := r.Form.Get("username")

		sessionToken, expiresAt, err := controller.HandleLogin(email, password, username)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    sessionToken,
			Expires:  expiresAt,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
