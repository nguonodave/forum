package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"real_time_forum/backend/db/controllers"
)

func LikedPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	// Get user ID from session
	cookie, err := controllers.GetCookie(r, "session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userId, _, _, err := controllers.GetUserDetailsFromSession(cookie)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Get liked posts
	posts, err := controllers.GetLikedPosts(userId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println(posts)
	// Return posts as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

type ReceivedRequest struct {
	Id       string `json:"id"`
	Reaction string `json:"reaction"`
}

func LikesHandler(w http.ResponseWriter, r *http.Request) {
	var (
		req    ReceivedRequest
		userId string
	)

	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(resBody, &req); err != nil {
		log.Println("Error unmarshalling JSON:", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	postId := req.Id
	reaction := req.Reaction
	likes := 0
	dislikes := -1

	cookie, err := controllers.GetCookie(r, "session_token")
	if err != nil {
		http.Error(w, "Unauthorized: Please SignIn/SignUp to Continue", http.StatusUnauthorized)
		return
	}

	if controllers.ValidSession(cookie) {
		userId, _, _, _ = controllers.GetUserDetailsFromSession(cookie)
	}

	check, err := controllers.CheckPreviousReaction(postId, userId)
	if err != nil {
		log.Println("Error checking previous reaction:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if check == "" { // User has never reacted before
		controllers.InsertReaction(userId, postId, reaction)
		if likes, err = controllers.IncreaseLikes(postId); err != nil {
			log.Println("Error increasing likes:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else if check == "like" {
		if reaction == "like" { // Unlike if already liked
			if likes, err = controllers.DecreaseLikes(postId); err != nil {
				log.Println("Error decreasing likes:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if err := controllers.DeleteReaction(postId, userId); err != nil {
				log.Println("Error deleting reaction:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
	} else if check == "dislike" {
		// Remove dislike and switch to like
		if dislikes, err = controllers.DecreaseDislikes(postId); err != nil {
			log.Println("Error decreasing dislikes:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if err := controllers.DeleteReaction(postId, userId); err != nil {
			log.Println("Error deleting reaction:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		controllers.InsertReaction(userId, postId, reaction)
		if likes, err = controllers.IncreaseLikes(postId); err != nil {
			log.Println("Error increasing likes:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct{
		Likes int 
		Dislikes int}{
		Likes: likes,
		Dislikes: dislikes,})
}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		req    ReceivedRequest
		userId string
	)

	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(resBody, &req); err != nil {
		log.Println("Error unmarshalling JSON:", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	postId := req.Id
	reaction := req.Reaction
	dislikes := 0
	likes := -1

	cookie, err := controllers.GetCookie(r, "session_token")
	if err != nil {
		http.Error(w, "Unauthorized: Please SignIn/SignUp to Continue", http.StatusUnauthorized)
		return
	}

	if controllers.ValidSession(cookie) {
		userId, _, _, _ = controllers.GetUserDetailsFromSession(cookie)
	}

	check, err := controllers.CheckPreviousReaction(postId, userId)
	if err != nil {
		log.Println("Error checking previous reaction:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if check == "" { // User has never reacted before
		controllers.InsertReaction(userId, postId, reaction)
		if dislikes, err = controllers.IncreaseDislikes(postId); err != nil {
			log.Println("Error increasing dislikes:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else if check == "dislike" {
		if reaction == "dislike" { // Unlike if already liked
			if dislikes, err = controllers.DecreaseDislikes(postId); err != nil {
				log.Println("Error decreasing likes:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if err := controllers.DeleteReaction(postId, userId); err != nil {
				log.Println("Error deleting reaction:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
	} else if check == "like" {
		// Remove like and switch to dislike
		if likes, err = controllers.DecreaseLikes(postId); err != nil {
			log.Println("Error decreasing dislikes:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if err := controllers.DeleteReaction(postId, userId); err != nil {
			log.Println("Error deleting reaction:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		controllers.InsertReaction(userId, postId, reaction)
		if dislikes, err = controllers.IncreaseDislikes(postId); err != nil {
			log.Println("Error increasing likes:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct{
		Likes int 
		Dislikes int}{
		Likes: likes,
		Dislikes: dislikes,})
}

func CommentsLikesHandler(w http.ResponseWriter, r *http.Request) {
	var userId string

	path := r.Referer()
	r.ParseForm()

	commentId := r.FormValue("commentid")
	reaction := r.FormValue("reaction")

	cookie, err := controllers.GetCookie(r, "session_token")
	if err != nil {
		// renderErrorPage(w, r, http.StatusUnauthorized, "Please SignIn/SignUp to Continue")
		return
	}

	if controllers.ValidSession(cookie) {
		userId, _, _, _ = controllers.GetUserDetailsFromSession(cookie)
	}

	check, err := controllers.CheckPreviousCommentReaction(commentId, userId)
	if err != nil {
		fmt.Println(err.Error())
		// renderErrorPage(w, r, http.StatusInternalServerError, "Internal servor error")
		return
	}

	if check == "" {
		controllers.InsertCommentReaction(userId, commentId, reaction)
		if reaction == "like" {
			if err := controllers.IncreaseCommentLikes(commentId); err != nil {
				// renderErrorPage(w, r, http.StatusInternalServerError, "Internal servor error")
				return
			}
		} else if reaction == "dislike" {
			if err := controllers.IncreaseCommentDislikes(commentId); err != nil {
				// renderErrorPage(w, r, http.StatusInternalServerError, "Internal servor error")
				return
			}
		}
	} else if check == "like" {
		if reaction == check {
			if err := controllers.DecreaseCommentLikes(commentId); err != nil {
				// renderErrorPage(w, r, http.StatusInternalServerError, "Internal servor error")
				return
			}
			if err := controllers.DeleteCommentReaction(commentId, userId); err != nil {
				// renderErrorPage(w, r, http.StatusInternalServerError, "Internal servor error")
				return
			}
		} else {
			if err := controllers.UpdateCommentReaction(reaction, commentId, userId); err != nil {
				// renderErrorPage(w, r, http.StatusInternalServerError, "Internal servor error")
				return
			}

			if err := controllers.IncreaseCommentDislikes(commentId); err != nil {
				// renderErrorPage(w, r, http.StatusInternalServerError, "Internal servor error")
				return
			}

			if err := controllers.DecreaseCommentLikes(commentId); err != nil {
				// renderErrorPage(w, r, http.StatusInternalServerError, "Internal servor error")
				return
			}
		}
	} else if check == "dislike" {
		if reaction == check {
			if err := controllers.DecreaseCommentDislikes(commentId); err != nil {
				// renderErrorPage(w, r, http.StatusInternalServerError, "Internal servor error")
				return
			}

			if err := controllers.DeleteCommentReaction(commentId, userId); err != nil {
				// renderErrorPage(w, r, http.StatusInternalServerError, "Internal servor error")
				return
			}

		} else {
			if err := controllers.UpdateCommentReaction(reaction, commentId, userId); err != nil {
				// renderErrorPage(w, r, http.StatusInternalServerError, "Internal servor error")
				return
			}

			if err := controllers.IncreaseCommentLikes(commentId); err != nil {
				// renderErrorPage(w, r, http.StatusInternalServerError, "Internal servor error")
				return
			}

			if err := controllers.DecreaseCommentDislikes(commentId); err != nil {
				// renderErrorPage(w, r, http.StatusInternalServerError, "Internal servor error")
				return
			}
		}
	}

	http.Redirect(w, r, path, http.StatusSeeOther)
}
