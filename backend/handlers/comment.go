package handlers

import (
	"encoding/json"
	"io"

	"net/http"
	"real_time_forum/backend/db/controllers"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type ReceivedComment struct {
	Comment string `json:"comment"`
	PostId  string `json:"postId"`
}

func HandleComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var comment ReceivedComment

	rBody, err := io.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": err.Error()})
		return
	}

	if err = json.Unmarshal(rBody, &comment); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": err.Error()})
		return
	}

	postId := comment.PostId
	commentContent := comment.Comment

	if postId == "" || strings.TrimSpace(commentContent) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Empty fields not allowed"})
		return
	}

	// Get user session
	cookie, err := controllers.GetCookie(r, "session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "unauthorized"})
		return
	}

	// Validate session
	if !controllers.ValidSession(cookie) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userId, userName, _, _ := controllers.GetUserDetailsFromSession(cookie)

	// Insert the comment into the database
	commentID, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Failed to generate comment ID"})
		return
	}

	CreatedAt := time.Now().Format("2006-01-02 15:04:05")

	comment_count := 0

	if comment_count, err = controllers.InsertCommentToDB(commentID.String(), postId, userId, userName, commentContent, CreatedAt); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Failed to save comment"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "comment saved successfully", 
		"commentId" : commentID.String(), 
		"postId" : postId, 
		"comment" : comment, 
		"time" : CreatedAt, 
		"comments_count" : comment_count})
	// fmt.Println(postId, commentContent, userId, userName)

}
