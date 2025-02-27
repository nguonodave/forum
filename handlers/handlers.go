package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
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
		ErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
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
			ErrorPage(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		sessionToken, expiresAt, err := controller.HandleLogin(data.Email, data.Username, data.Password)
		if err != nil {
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
		ErrorPage(w, "Method not allowed", http.StatusMethodNotAllowed)
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
			ErrorPage(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := controller.HandleRegister(data.Username, data.Email, data.Password); err != nil {
			jsonResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		jsonResponse(w, http.StatusOK, "Registration successful")

	default:
		ErrorPage(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleVoteRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		ErrorPage(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// since this handler is wrapped with the Validate Session handler, which comes with user id and username already from r.Context no need to get user id again
	ctxUserID, ok := r.Context().Value("userId").(string)
	if !ok || ctxUserID == "" {
		ErrorPage(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	// parse the user id to validate if its a valid uuid
	userId, err := uuid.Parse(ctxUserID)
	if err != nil {
		ErrorPage(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	if userId.String() == "" {
		userId, err = helperfunc.GetUserIDFromSession(r)
		if err != nil {
			ErrorPage(w, "User not authenticated", http.StatusUnauthorized)
			return
		}
	}

	// Parse request body
	var requestBody struct {
		PostID    string `json:"postId,omitempty"`
		CommentID string `json:"commentId,omitempty"`
		Type      string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		ErrorPage(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate vote type
	if requestBody.Type != "like" && requestBody.Type != "dislike" {
		ErrorPage(w, "Invalid vote type", http.StatusBadRequest)
		return
	}

	// Prepare vote request
	voteReq := &model.VoteRequest{
		UserID: userId,
		Type:   requestBody.Type,
	}

	if requestBody.PostID != "" {
		postID, err := uuid.Parse(requestBody.PostID)
		if err != nil {
			ErrorPage(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		voteReq.PostID = &postID
	} else if requestBody.CommentID != "" {
		commentID, err := uuid.Parse(requestBody.CommentID)
		if err != nil {
			ErrorPage(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}
		voteReq.CommentID = &commentID
	} else {
		ErrorPage(w, "Must provide either postId or commentId", http.StatusBadRequest)
		return
	}

	// Process vote
	response, err := model.HandleVote(voteReq)

	if err != nil {
		ErrorPage(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		ErrorPage(w, "no session found", http.StatusUnauthorized)
		return
	}

	sessionToken := cookie.Value

	query := "DELETE FROM sessions WHERE token = ?"
	result, err := database.Db.Exec(query, sessionToken)
	if err != nil {
		ErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ErrorPage(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		ErrorPage(w, "session not found", http.StatusInternalServerError)
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

func AddCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			ErrorPage(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		//	decode JSON request body
		var commentReq model.CommentRequest

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&commentReq)
		if err != nil {
			ErrorPage(w, "Invalid request JSON", http.StatusBadRequest)
			return
		}
		userId := r.Context().Value("userId").(string)
		commentReq.UserID = userId

		// validate input
		if commentReq.PostID == "" || commentReq.UserID == "" || strings.TrimSpace(commentReq.Content) == "" {
			ErrorPage(w, "missing required fields", http.StatusBadRequest)
			return
		}

		commentID := uuid.New().String()
		createdAt := time.Now().Format(time.ANSIC)
		query := `INSERT INTO comments (id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)`
		_, err = db.Exec(query, commentID, commentReq.PostID, commentReq.UserID, commentReq.Content, createdAt)
		if err != nil {
			ErrorPage(w, "failed to insert comment", http.StatusInternalServerError)
			return
		}

		// return this response
		comment := Comment{
			ID:      commentID,
			PostID:  commentReq.PostID,
			UserID:  commentReq.UserID,
			Content: commentReq.Content,
		}

		// send JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(comment)
	}
}

type Comment struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
