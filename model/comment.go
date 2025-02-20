package model

type CommentRequest struct {
	PostID  string `json:"post_id"`
	UserID  string `json:"user_id,omitempty"`
	Content string `json:"content"`
}

type CommentResponse struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type Comment struct {
	Id        string
	PostId    string
	UserId    string
	Content   string
	Likes     int
	Dislikes  int
	CreatedAt string
}
