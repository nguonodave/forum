package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"forum/controller"
	"forum/database"
	"forum/helperfunc"
	"forum/model"

	"github.com/google/uuid"
)

// templatesDir refers to the filepath of the directory containing the application's templates
var templatesDir = "view"

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

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		Templates.ExecuteTemplate(w, "auth.html", nil)

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

		sessionToken, expiresAt, err := controller.HandleLogin(data.Email, data.Username, data.Password)
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
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		Templates.ExecuteTemplate(w, "auth.html", nil)

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
		if err := controller.HandleRegister(data.Username, data.Email, data.Password); err != nil {
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

func HandleVoteRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("/api/vote has been hit...")
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	// since this handler is wrapped with the Validate Session handler, which comes with user id and username already from r.Context no need to get user id again
	ctxUserID, ok := r.Context().Value("userId").(string)
	fmt.Println("user id from context and ok", ctxUserID, ok)
	if !ok || ctxUserID == "" {
		fmt.Println(3)
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	fmt.Println(4)
	// parse the user id to validate if its a valid uuid
	userId, err := uuid.Parse(ctxUserID)
	fmt.Println(5)
	if err != nil {
		fmt.Println(6)
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	if userId.String() == "" {
		fmt.Println(7)
		userId, err = helperfunc.GetUserIDFromSession(r)
		fmt.Println(8)
		if err != nil {
			fmt.Println(9)
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}
	}

	// Parse request body
	var requestBody struct {
		PostID    string `json:"postId,omitempty"`
		CommentID string `json:"commentId,omitempty"`
		Type      string `json:"type"`
	}
	fmt.Println(10)

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Println(11)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate vote type
	if requestBody.Type != "like" && requestBody.Type != "dislike" {
		fmt.Println(12)
		fmt.Println("CANNOT POST DUE TO DB ERROR....handlers.go line 165")
		http.Error(w, "Invalid vote type", http.StatusBadRequest)
		return
	}

	// Prepare vote request
	voteReq := &model.VoteRequest{
		UserID: userId,
		Type:   requestBody.Type,
	}

	// Set PostID or CommentID based on request
	if requestBody.PostID != "" {
		postID, err := uuid.Parse(requestBody.PostID)
		fmt.Println(12)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		voteReq.PostID = &postID
	} else if requestBody.CommentID != "" {
		fmt.Println(13)
		commentID, err := uuid.Parse(requestBody.CommentID)
		if err != nil {
			fmt.Println(14)
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}
		voteReq.CommentID = &commentID
	} else {
		fmt.Println(15)
		http.Error(w, "Must provide either postId or commentId", http.StatusBadRequest)
		return
	}

	// Process vote
	response, err := model.HandleVote(voteReq)
	fmt.Println(err, 16)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(17)
	// Send JSON response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse(w, http.StatusMethodNotAllowed, http.StatusText(405))
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "no session found", http.StatusUnauthorized)
		return
	}

	sessionToken := cookie.Value

	query := "DELETE FROM sessions WHERE token = ?"
	result, err := database.Db.Exec(query, sessionToken)
	if err != nil {
		fmt.Printf("111%v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("112%v", err)
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
	jsonResponse(w, http.StatusOK, "Logout successful")
}
