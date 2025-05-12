package main

import (
	"log"
	"net/http"
	"real_time_forum/backend/db/controllers"
	"real_time_forum/backend/handlers"
	"real_time_forum/backend/socket"
)

func init() {
	controllers.InitDB()
	controllers.CreateTables()
}
func main() {
	file := "frontend/templates/index.html"
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer((http.Dir("frontend/static")))))

	http.HandleFunc("/login", handlers.HandleLogin)
	http.HandleFunc("/signup", handlers.SignupHandler)
	http.HandleFunc("/home", handlers.LoadAll)
	http.HandleFunc("/createPost", handlers.PostHandler)
	http.HandleFunc("/likePost", handlers.LikesHandler)
	http.HandleFunc("/dislikePost", handlers.DislikeHandler)
	http.HandleFunc("/createComment", handlers.HandleComment)
	http.HandleFunc("/privateMessage", handlers.HandlePrivateMessage)
	http.HandleFunc("/getusermessages", handlers.GetUserMessages)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/ws", socket.HandleRequest)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, file)
	})
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
