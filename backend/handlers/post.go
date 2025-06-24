package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"real_time_forum/backend/db/controllers"

	"github.com/gofrs/uuid"
)
type Received struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Categories string `json:"categories"`
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}


	var received Received

	// Read the request body
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "Failed to read request body"}`, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Unmarshal JSON
	if err := json.Unmarshal(requestBody, &received); err != nil {
		log.Println("Error unmarshaling post data:", err)
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}
	
	// Get user session
	cookie, err := controllers.GetCookie(r, "session_token")
	if err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}
	
	// Validate session
	if !controllers.ValidSession(cookie) {
		w.Header().Set("Conten-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message" : "session_expired please login"})
		return
	}

	userId, username, _, _ := controllers.GetUserDetailsFromSession(cookie)
	// log.Println("User details ID:", userId, username, identifier)
	
	Handlepost(w, r, userId, username, received)

}


func Handlepost(w http.ResponseWriter, r *http.Request, userId, username string, received Received) {
	var (
		id        uuid.UUID
		title     string
		content   string
		createdAt string
		categories  string
		err       error
	)

	id, err = uuid.NewV4()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	title = received.Title
	content = received.Content
	createdAt = time.Now().Format("2006-01-02 15:04:05")
	categories = received.Categories
	
	
	// log.Println("Received Data:", id, userId, username, title, content, createdAt, categories)
	if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
		http.Error(w, "Title and content cannot be empty", http.StatusBadRequest)
		return
	}

	err = controllers.AddNewPostToDB(id, userId, username, title, content, createdAt, categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Post created successfully", "postId" : id})
}
