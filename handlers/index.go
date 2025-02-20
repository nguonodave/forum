package handlers

import (
	"database/sql"
	"fmt"
	"forum/model"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"forum/controller"

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

	queriedCategory := r.URL.Query().Get("category")

	if r.Method == http.MethodPost {
		userLoggedIn, username, userId := pkg.UserLoggedIn(r)
		if !userLoggedIn {
			http.Error(w, "session expired", http.StatusUnauthorized)
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

		for _, category := range categories {
			_, insertCategoriesErr := database.Db.Exec(`INSERT INTO post_categories (post_id, category) VALUES (?, ?)`, postId, category)
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

	// rendering posts to template
	type Post struct {
		Username     string
		Id           string
		Title        string
		Content      string
		ImagePath    string
		CreatedTime  string
		Likes        int
		Dislikes     int
		Categories   []string
		CommentCount int
		Comments     []model.Comment
	}

	var posts []Post

	postsQuery := `
	SELECT DISTINCT p.username , p.id, p.title, p.content, p.image_url, p.created_at
	FROM posts p
	LEFT JOIN post_categories pc ON p.id = pc.post_id
	WHERE ? = '' OR pc.category = ?
	ORDER BY p.created_at DESC
	`

	rows, err := database.Db.Query(postsQuery, queriedCategory, queriedCategory)
	if err != nil {
		log.Printf("Error fetching posts: %v\n", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		var imagePath sql.NullString // use sql.NullString to handle NULL values
		var formatTime time.Time

		err := rows.Scan(&post.Username, &post.Id, &post.Title, &post.Content, &imagePath, &formatTime)
		if err != nil {
			log.Printf("Error scanning post: %v\n", err)
			continue
		}

		post.CreatedTime = formatTime.Format(time.ANSIC)

		// if imagePath is valid, assign it to the post
		if imagePath.Valid {
			post.ImagePath = imagePath.String
		}

		// get all the categories of the post
		postCategories, postCategoriesErr := PostCategories(post.Id)
		if postCategoriesErr != nil {
			http.Error(w, "Failed to fetch post categories", http.StatusInternalServerError)
			return
		}

		post.Categories = append(post.Categories, postCategories...)

		posts = append(posts, post)
	}

	for i := range posts { // Use index-based iteration to modify slice elements
		err := controller.GetLikesDislikesForPost(database.Db, posts[i].Id, &posts[i].Likes, &posts[i].Dislikes)
		if err != nil {
			log.Printf("Error fetching post likes: %v\n", err)
		}

		// comments is a slice of comments for a certain post with id post[i].Id
		comments, err := controller.FetchCommentsForPost(database.Db, posts[i].Id)
		if err != nil {
			log.Printf("Error fetching post comments: %v\n", err)
		}
		posts[i].CommentCount = len(comments)
		posts[i].Comments = comments
	}

	categories, getCategoriesErr := pkg.GetCategories(w)
	if getCategoriesErr != nil {
		return
	}

	TemplateError := func(message string, err error) {
		http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
		log.Printf("%s: %v", message, err)
	}

	userLoggedIn, username, _ := pkg.UserLoggedIn(r)

	execTemplateErr := Templates.ExecuteTemplate(w, "base.html", map[string]interface{}{
		"Posts":        posts,
		"UserLoggedIn": userLoggedIn,
		"Categories":   categories,
		"Username":     username,
	})
	if execTemplateErr != nil {
		TemplateError("error executing template", execTemplateErr)
		return
	}
}

func PostCategories(postId string) ([]string, error) {
	categories := []string{}

	// get all the categories of the post
	postCategoriesRows, postCategoriesRowsErr := database.Db.Query(`SELECT category FROM post_categories WHERE post_id = ?`, postId)
	if postCategoriesRowsErr != nil {
		log.Printf("Error getting categories for post: %v\n", postCategoriesRowsErr)
		return nil, postCategoriesRowsErr
	}
	defer postCategoriesRows.Close()

	for postCategoriesRows.Next() {
		var category string
		postCategoriesScanErr := postCategoriesRows.Scan(&category)
		if postCategoriesScanErr != nil {
			log.Printf("Error scanning post categories: %v\n", postCategoriesScanErr)
			return nil, postCategoriesScanErr
		}
		categories = append(categories, category)
	}

	return categories, nil
}
