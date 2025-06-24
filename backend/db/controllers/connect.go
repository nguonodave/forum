package controllers

import (
	"database/sql"
	"log"

	"real_time_forum/backend/db/models"
	"real_time_forum/backend/utils"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB  *sql.DB
	err error
)

func InitDB() {
	DB, err = sql.Open("sqlite3", "./backend/db/forum.db")
	if err != nil {
		log.Fatal("DB init error", err.Error())
	}

	// defer DB.Close()
	// Test the database connection
	if err := DB.Ping(); err != nil {
		log.Fatal("Ping error:", err.Error())
	}
}

func CreateTables() error {
	for _, stmt := range Tables {
		_, err := DB.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func All() ([]models.Post, error) {
	// Prepare the SQL query
	stm, err := DB.Prepare(`
	SELECT
		posts.id AS postId,
		posts.userid AS userID,
		posts.username AS postsUser,
		posts.title AS postTitle,
		posts.content AS postContent,
		posts.createdAt AS creationDate,
		posts.category AS categories,
		posts.likes AS likes,
		posts.dislikes AS dislikes,
		comments.id AS commentId,
		comments.userid AS commenterId,
		comments.username AS commenter,
		comments.postid AS postIdForComment,
		comments.content AS commentText,
		comments.createdAt AS commentCreationDate,
		comments.likes AS commentLikes,
		comments.dislikes AS commentDislikes
	FROM
		posts
	LEFT JOIN
		comments
	ON
		posts.id = comments.postid
	`)
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	rows, err := stm.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	postsMap := make(map[string]*models.Post)

	for rows.Next() {
		var post models.Post
		var comment models.Comment
		var commentID, commenterID, commentContent, commentCreatedAt , commenter sql.NullString
		var commentPostID sql.NullString
		var commentLikes, commentDislikes sql.NullInt64

		// Scan the row
		if err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Name,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Category,
			&post.Likes,
			&post.Dislikes,
			&commentID,
			&commenterID,
			&commenter,
			&commentPostID,
			&commentContent,
			&commentCreatedAt,
			&commentLikes,
			&commentDislikes,
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Populate comment if it exists
		if commentID.Valid {
			comment.ID = commentID.String
			comment.UserID = commenterID.String
			comment.UserName = commenter.String
			comment.PostId = commentPostID.String
			comment.Content = commentContent.String
			comment.CreatedAt = commentCreatedAt.String
			comment.Likes = int(commentLikes.Int64)
			comment.Dislikes = int(commentDislikes.Int64)
		}

		// Add to posts map
		if existingPost, found := postsMap[post.ID]; found {
			if commentID.Valid {
				existingPost.Comments = append(existingPost.Comments, comment)
			}
		} else {
			post.Comments = []models.Comment{}
			if commentID.Valid {
				post.Comments = append(post.Comments, comment)
			}
			postsMap[post.ID] = &post
		}
		// fmt.Println(len(post.Comments))
	}

	// Collect posts into a slice
	posts := utils.SortPostsByLatest(postsMap)

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// log.Printf("Loaded %d posts", len(posts))
	return posts, nil
}
