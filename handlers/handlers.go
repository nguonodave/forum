package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"forum/controller"
	"forum/helperfunc"
	"forum/model"

	"github.com/google/uuid"
)

// templatesDir refers to the filepath of the directory containing the application's templates
var templatesDir = "view"

// renderTemplate is a helper function to render HTML templates
func renderTemplate(w http.ResponseWriter, templateFile string, data interface{}) {
	templatePath := filepath.Join(templatesDir, templateFile)
	temp, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
		log.Printf("Error parsing template %s: %v", templateFile, err)
		return
	}

	err = temp.Execute(w, data)
	if err != nil {
		fmt.Println("rrr")
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
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func Login(DBase *model.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			renderTemplate(w, "auth/auth.html", nil)

		case http.MethodPost:
			var data struct {
				Email    string `json:"email,omitempty"`
				Username string `json:"username,omitempty"`
				Password string `json:"password"`
			}

			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				log.Printf("error decoding login data: %v", err)
				jsonResponse(w, http.StatusBadRequest, "Invalid JSON")
				return
			}

			sessionToken, expiresAt, err := controller.HandleLogin(DBase, data.Email, data.Username, data.Password)
			if err != nil {
				log.Printf("%v", err)
				jsonResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			// Set session cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    sessionToken,
				Expires:  expiresAt,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})

			jsonResponse(w, http.StatusOK, "Login successful")
			http.Redirect(w, r, "/", http.StatusFound)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func Register(DBase *model.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			renderTemplate(w, "auth/auth.html", nil)

		case http.MethodPost:
			var data struct {
				Username string `json:"username"`
				Email    string `json:"email"`
				Password string `json:"password"`
			}

			// Decode JSON request
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				log.Printf("%v", err)
				jsonResponse(w, http.StatusBadRequest, "Invalid JSON")
				return
			}

			// Register user
			if err := controller.HandleRegister(DBase, data.Username, data.Email, data.Password); err != nil {
				fmt.Printf("%v", err)
				jsonResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			fmt.Println("Registration successful")
			jsonResponse(w, http.StatusOK, "Registration successful")

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func CategoriesHandler(db *model.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := model.GetCategories(db.Db)
		if err != nil {
			http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles(".html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, categories)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
		}
	}
}

func HandleVoteRequest(db *model.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set response header
		w.Header().Set("Content-Type", "application/json")

		// Only allow POST method
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get user ID from session
		userID, err := helperfunc.GetUserIDFromSession(r)
		if err != nil {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		// Parse request body
		var requestBody struct {
			PostID    string `json:"postId,omitempty"`
			CommentID string `json:"commentId,omitempty"`
			Type      string `json:"type"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate vote type
		if requestBody.Type != "like" && requestBody.Type != "dislike" {
			http.Error(w, "Invalid vote type", http.StatusBadRequest)
			return
		}

		// Prepare vote request
		voteReq := &model.VoteRequest{
			UserID: userID,
			Type:   requestBody.Type,
		}

		// Set PostID or CommentID based on request
		if requestBody.PostID != "" {
			postID, err := uuid.Parse(requestBody.PostID)
			if err != nil {
				http.Error(w, "Invalid post ID", http.StatusBadRequest)
				return
			}
			voteReq.PostID = &postID
		} else if requestBody.CommentID != "" {
			commentID, err := uuid.Parse(requestBody.CommentID)
			if err != nil {
				http.Error(w, "Invalid comment ID", http.StatusBadRequest)
				return
			}
			voteReq.CommentID = &commentID
		} else {
			http.Error(w, "Must provide either postId or commentId", http.StatusBadRequest)
			return
		}

		// Process vote
		response, err := model.HandleVote(db, voteReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send response
		json.NewEncoder(w).Encode(response)
	}
}

func Logout(db *model.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			jsonResponse(w, http.StatusMethodNotAllowed, "not allowed")
			return
		}
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "no session found", http.StatusUnauthorized)
			return
		}

		sessionToken := cookie.Value

		query := "DELETE FROM sessions WHERE token = ?"
		result, err := db.Db.Exec(query, sessionToken)
		if err != nil {
			fmt.Printf("%v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			fmt.Printf("%v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "session not found", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
