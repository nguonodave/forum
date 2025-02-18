package handlers

import (
	"database/sql"
	"fmt"
	"forum/controller"
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

	queriedCategoryId := r.URL.Query().Get("category")

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
			http.Error(w, "An unexpected error occurred. Please try again later", http.StatusInternalServerError)
			return
		}
		// get the username from users table
		var username string
		usernameGetError := database.Db.QueryRow("SELECT username FROM users WHERE id = ?", userId).Scan(&username)
		if usernameGetError != nil {
			log.Printf("Error fetching username from database: %v\n", usernameGetError)
			http.Error(w, "An unexpected error occurred. Please try again later", http.StatusInternalServerError)
			return
		}
		fmt.Println("user id", userId)

		// 20 MB limit
		maxSizeErr := r.ParseMultipartForm(20 << 20)
		if maxSizeErr != nil {
			http.Error(w, "Unable to parse form. Make sure image size is not more than 20mb", http.StatusBadRequest)
			return
		}

		// get the submitted form content details
		title := r.FormValue("title")
		content := r.FormValue("content")
		categories := r.Form["categories"]
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

		_, insertPostErr := database.Db.Exec(`INSERT INTO posts (username , id, user_id, title, content, image_url, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`, username, postId, userId, title, content, imagePath, createdTime)
		if insertPostErr != nil {
			log.Printf("Error inserting post to database: %v\n", insertPostErr)
			http.Error(w, "Failed to add post", http.StatusInternalServerError)
			return
		}

		for _, categoryId := range categories {
			_, insertCategoriesErr := database.Db.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, postId, categoryId)
			if insertCategoriesErr != nil {
				log.Printf("Error inserting selected categories to database: %v\n", insertCategoriesErr)
				http.Error(w, "Failed to add categories", http.StatusInternalServerError)
				return
			}
		}

		// redirect to the home page after creating post
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	type Post struct {
		Username    string
		Id          string
		Title       string
		Content     string
		ImagePath   string
		CreatedTime string
		Likes       int
		Dislikes    int
	}

	var posts []Post

	postsQuery := `
	SELECT DISTINCT p.username , p.id, p.title, p.content, p.image_url, p.created_at
	FROM posts p
	LEFT JOIN post_categories pc ON p.id = pc.post_id
	WHERE ? = '' OR pc.category_id = ?
	ORDER BY p.created_at DESC
	`

	rows, err := database.Db.Query(postsQuery, queriedCategoryId, queriedCategoryId)
	if err != nil {
		log.Printf("Error fetching posts: %v\n", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var post Post

	for rows.Next() {

		var imagePath sql.NullString
		var formatTime time.Time
		err := rows.Scan(&post.Username, &post.Id, &post.Title, &post.Content, &imagePath, &formatTime)
		if err != nil {
			log.Printf("Error scanning post: %v\n", err)
			continue
		}

		post.CreatedTime = formatTime.Format(time.ANSIC)

		if imagePath.Valid {
			post.ImagePath = imagePath.String
		}

		posts = append(posts, post)
	}
	for i := range posts { // Use index-based iteration to modify slice elements
		fmt.Println("before", posts[i].Dislikes, posts[i].Likes)
		err := controller.GetLikesDislikesForPost(database.Db, posts[i].Id, &posts[i].Likes, &posts[i].Dislikes)
		if err != nil {
			log.Printf("Error fetching post likes: %v\n", err)
		}
		fmt.Println("AFTER", posts[i].Dislikes, posts[i].Likes)
	}

	// fetch all categories to render to the create post form
	categRows, categQueryErr := database.Db.Query(`SELECT id, name FROM categories`)
	if categQueryErr != nil {
		log.Printf("Error fetching categories: %v\n", categQueryErr)
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}
	defer categRows.Close()

	var categories []struct {
		Id   string
		Name string
	}
	for categRows.Next() {
		var category struct {
			Id   string
			Name string
		}
		err := categRows.Scan(&category.Id, &category.Name)
		if err != nil {
			log.Printf("Error scanning category: %v\n", err)
			continue
		}
		categories = append(categories, category)
	}

	TemplateError := func(message string, err error) {
		http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
		log.Printf("%s: %v", message, err)
	}

	fmt.Printf("????%+v\n", posts)
	isLoggedIn, username := pkg.UserLoggedIn(r)
	fmt.Println(isLoggedIn, username)
	execTemplateErr := Templates.ExecuteTemplate(w, "base.html", map[string]interface{}{
		"Posts":        posts,
		"UserLoggedIn": isLoggedIn,
		"Categories":   categories,
		"Username":     username,
	})

	if execTemplateErr != nil {
		TemplateError("error executing template", execTemplateErr)
		return
	}
}
