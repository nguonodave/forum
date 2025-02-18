package posts

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"forum/database"
	"forum/handlers"
	"forum/pkg"
)

func CategoryPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handlers.ErrorPage(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	queriedCategory := strings.TrimPrefix(r.URL.Path, "/posts/")

	// rendering posts to template
	type Post struct {
		Id          string
		Title       string
		Content     string
		ImagePath   string
		CreatedTime time.Time
	}

	var posts []Post

	postsQuery := `
	SELECT DISTINCT p.id, p.title, p.content, p.image_url, p.created_at
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

	// populate the post struct
	for rows.Next() {
		var post Post
		var imagePath sql.NullString // use sql.NullString to handle NULL values
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &imagePath, &post.CreatedTime)
		if err != nil {
			log.Printf("Error scanning post: %v\n", err)
			continue
		}

		// if imagePath is valid, assign it to the post
		if imagePath.Valid {
			post.ImagePath = "../" + imagePath.String
		}

		posts = append(posts, post)
	}

	categories, getCategoriesErr := pkg.GetCategories(w)
	if getCategoriesErr != nil {
		return
	}

	TemplateError := func(message string, err error) {
		http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
		log.Printf("%s: %v", message, err)
	}

	execTemplateErr := handlers.Templates.ExecuteTemplate(w, "categoryposts.html", map[string]interface{}{
		"Posts":        posts,
		"UserLoggedIn": pkg.UserLoggedIn(r),
		"Categories":   categories,
	})
	if execTemplateErr != nil {
		TemplateError("error executing template", execTemplateErr)
		return
	}
}
