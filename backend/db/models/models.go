package models

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Age       int    `json:"age"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Post struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	CreatedAt   string    `json:"created_at"`
	Category    string    `json:"category"`
	Comments    []Comment `json:"comments"`
	CommentsNum int       `json:"comments_num"`
	Likes       int       `json:"likes"`
	Dislikes    int       `json:"dislikes"`
}
type Comment struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	PostId    string `json:"post_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
}

type Session struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ExpiresAt string `json:"expires_at"`
}
