package controllers

import ("real_time_forum/backend/db/models"

)

func GetLikedPosts(userId string) ([]models.Post, error) {
query := `
SELECT p.id, p.userid, p.username, p.title, p.content, p.createdat, p.category FROM posts p
INNER JOIN likes l ON p.id = l.postid
WHERE l.userid = ? AND l.reaction = 'like'
`
rows, err := DB.Query(query, userId)

if err!= nil{
	return nil, err
}
defer rows.Close()

var posts []models.Post

for rows.Next(){
	var post models.Post
	err := rows.Scan(
		&post.ID,
		&post.UserID,
		&post.Name,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.Category,
	)
	if err != nil {
		return nil, err
	}
	posts = append(posts, post)

}

return posts, nil
}
