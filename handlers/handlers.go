package handlers

import (
	"encoding/json"
	"flag"
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

var port = flag.Int("P", 8080, "port to listen on")

// templatesDir refers to the filepath of the directory containing the application's templates
var templatesDir = "view"

// renderTemplate is a helper function to render HTML templates
func renderTemplate(w http.ResponseWriter, templateFile string, data interface{}) {
	templatePath := filepath.Join(templatesDir, templateFile)
	temp, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error parsing template %s: %v", templateFile, err)
		return
	}

	err = temp.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template %s: %v", templateFile, err)
		return
	}
}

// jsonResponse is a helper function to send JSON responses
func jsonResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]string{"message": message}
	json.NewEncoder(w).Encode(response)
}

// Login handles both GET and POST methods
func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderTemplate(w, "auth/login.html", nil)

	case http.MethodPost:
		var data struct {
			Email    string `json:"email"`
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			jsonResponse(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		fmt.Printf("Received login data: Email=%s, Username=%s, Password=%s\n", data.Email, data.Username, data.Password)

		sessionToken, expiresAt, err := controller.HandleLogin(data.Email, data.Password)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, err.Error())
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

		jsonResponse(w, http.StatusOK, "Login successful")

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Register handles /register endpoint for user registration
func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderTemplate(w, "auth/login.html", nil)

	case http.MethodPost:
		var data struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			fmt.Println(err, 1)
			jsonResponse(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		fmt.Printf("Received registration data: Username=%s, Email=%s, Password=%s\n", data.Username, data.Email, data.Password)

		if err := controller.HandleRegister(data.Username, data.Email, data.Password); err != nil {
			fmt.Println(err, 2)
			jsonResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Println("Registration successful")
		jsonResponse(w, http.StatusOK, "Registration successful")

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
