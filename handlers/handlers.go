package handlers

import (
	"encoding/json"
	"fmt"
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
		fmt.Printf("body  %s\n", r.Body)
		var data struct {
			Email    string `json:"email"`
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		fmt.Println("Received data:", data.Email, data.Username, data.Password)
		sessionToken, expiresAt, err := controller.HandleLogin(data.Email, data.Password)
		fmt.Println("session", sessionToken, "expires at", expiresAt)
		println()
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Login successful"}`))
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
		var data struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		fmt.Println("Received data:", data.Username, data.Email, data.Password)
		err = controller.HandleRegister(data.Username, data.Email, data.Password)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Registration successful"}`))
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
