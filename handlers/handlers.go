package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"forum/controller"
	"forum/database"
	"forum/model"
)

// templatesDir refers to the filepath of the directory containing the application's templates
var templatesDir = "view"

// Login handles both GET and POST methods, if method is GET it renders the page
// if method is POST it gets the values from the form and internally checks if details exist in the database
func Login(w http.ResponseWriter, r *http.Request) {
	templateFile := "auth/login.html"
	fmt.Println("resp", *r)
	if r.Method == "GET" {
		println(2)
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
		fmt.Println(email, password)
		sessionToken, expiresAt, err := controller.VerifyLogin(email, password)
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

// Register handles /register endpoint for registering
func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		username := r.Form.Get("username")
		password := r.Form.Get("password")
		email := r.Form.Get("email")

		err = controller.RegisterUser(username, email, password)
		if err != nil {
			http.Error(w, "error during registration", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)

	case http.MethodGet:
		templateFile := "auth/login.html"

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
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func GetPaginatedPostsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.InitializeDB()
	if err != nil {
		http.Error(w, "Failed to initialize database", http.StatusInternalServerError)
	}

	// Get `page` and `limit` from query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10 // Default limit
	}

	offset := (page - 1) * limit

	posts, err := model.GetPaginatedPosts(db, limit, offset)
	if err != nil {
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
