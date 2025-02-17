package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"forum/database"
	"forum/pkg"

	"github.com/google/uuid"
)

// Index handler designed for the application's index page
func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorPage(w, "Page not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		cookie, cookieErr := r.Cookie("session")
		if cookieErr != nil {
			log.Printf("Error getting session's cookie from request: %v\n", cookieErr)
			http.Error(w, "An unexpected error occured. Please try again later", http.StatusInternalServerError)
			return
		}

		// get the user id to populate in post details
		var userId string
		sessionQueryErr := database.Db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", cookie.Value).Scan(&userId)
		if sessionQueryErr != nil {
			log.Printf("Error fetching session's user id from database: %v\n", sessionQueryErr)
			http.Error(w, "An unexpected error occured. Please try again later", http.StatusInternalServerError)
			return
		}

		// 20 MB limit
		maxSizeErr := r.ParseMultipartForm(20 << 20)
		if maxSizeErr != nil {
			http.Error(w, "Unable to parse form. Make sure image size is not more than 20mb", http.StatusBadRequest)
			return
		}

		// get the submitted form content details
		title := r.FormValue("title")
		content := r.FormValue("content")
		image, header, formFileErr := r.FormFile("image")

		var imagePath string
		if formFileErr == nil {
			defer image.Close()

			path := "./static/images/uploads/posts"

			// create the dir that will store the file
			mkdirErr := os.MkdirAll(path, os.ModePerm)
			if mkdirErr != nil {
				log.Printf("Error creating posts uploads directory: %v\n", mkdirErr)
				http.Error(w, "Failed to save image", http.StatusInternalServerError)
				return
			}

			// unique filename for the image using dates and file name
			imageName := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)
			imagePath = filepath.Join(path, imageName)

			// create the file path
			dst, err := os.Create(imagePath)
			if err != nil {
				log.Printf("Error creating file: %v\n", err)
				http.Error(w, "Failed to save image", http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			// copy the image from the form to the created path
			_, err = io.Copy(dst, image)
			if err != nil {
				log.Printf("Error saving file: %v\n", err)
				http.Error(w, "Failed to save image", http.StatusInternalServerError)
				return
			}
		}

		postId := uuid.New().String()
		createdTime := time.Now()

		_, insertPostErr := database.Db.Exec(`INSERT INTO posts (id, user_id, title, content, image_url, created_at) VALUES (?, ?, ?, ?, ?, ?)`, postId, userId, title, content, imagePath, createdTime)
		if insertPostErr != nil {
			log.Printf("Error inserting post to database: %v\n", insertPostErr)
			http.Error(w, "Failed to add post", http.StatusInternalServerError)
			return
		}
	}

	type Post struct {
		AuthorId    string
		Title       string
		Content     string
		ImagePath   string
		CreatedTime time.Time
	}

	var posts []Post

	rows, err := database.Db.Query(`
		SELECT user_id, title, content, image_url, created_at
		FROM posts
		ORDER BY created_at DESC
	`)
	if err != nil {
		log.Printf("Error fetching posts: %v\n", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		var imagePath sql.NullString // Use sql.NullString to handle NULL values
		err := rows.Scan(&post.AuthorId, &post.Title, &post.Content, &imagePath, &post.CreatedTime)
		if err != nil {
			log.Printf("Error scanning post: %v\n", err)
			continue
		}

		// If imagePath is valid, assign it to the post
		if imagePath.Valid {
			post.ImagePath = imagePath.String
		}
		
		posts = append(posts, post)
	}

	TemplateError := func(message string, err error) {
		http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
		log.Printf("%s: %v", message, err)
	}

	execTemplateErr := Templates.ExecuteTemplate(w, "base.html", map[string]interface{}{
		"Posts": posts,
		"UserLoggedIn": pkg.UserLoggedIn(r),
	})
	if execTemplateErr != nil {
		TemplateError("error executing template", execTemplateErr)
		return
	}
}
