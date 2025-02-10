package main

import (
	"flag"
	"fmt"
	"forum/controller"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	"forum/database"
	"forum/fileio"
	"forum/handlers"
)

var (
	port = flag.Int("P", 8080, "port to listen on")
	open = flag.Bool("O", false, "open server index page in the default browser")
)

func main() {
	// parse the defined command-line flags
	flag.Parse()

	// configure file logging to temporary application logger file
	{
		logFilePath := path.Join(os.TempDir(), fmt.Sprintf("%d-forum-logger.log", os.Getpid()))
		logger, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			log.Printf("failed to setup file logging: logging to stderr instead: %v\n", err)
		} else {
			log.Printf("saving logs to: %s\n", logFilePath)
		}
		log.SetOutput(logger)
		defer fileio.Close(logger)
	}

	// Initialize database
	db, err := database.InitializeDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Db.Close() // ✅ Correctly defer closing the database connection
	fmt.Println("Database initialized successfully!")

	// Register handlers with database dependency
	http.HandleFunc("/", controller.ValidateSession(db, handlers.Index))
	http.HandleFunc("/login", handlers.Login(db))
	http.HandleFunc("/register", handlers.Register(db))
	http.HandleFunc("/api/posts", handlers.GetPaginatedPostsHandler(db))

	// Browsers ping for the /favicon.ico icon, redirect to the respective static file
	http.Handle("/favicon.ico", http.RedirectHandler("/static/svg/favicon.svg", http.StatusFound))

	// Serve static files from the static dir, but ensure not to expose the directory entries
	staticDirFileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		reqPath := filepath.Clean(r.URL.Path)
		switch reqPath {
		case "/static", "/static/css", "/static/js", "/static/svg":
			handlers.ErrorPage(w, "Bad Request", http.StatusBadRequest)
			return
		}
		staticDirFileServer.ServeHTTP(w, r)
	})

	// Start server
	servePort := fmt.Sprintf(":%d", *port)
	url := fmt.Sprintf("http://localhost%s\n", servePort)
	fmt.Printf("Server running at %s\n", url)
	if *open {
		openBrowser(url)
	}
	log.Fatal(http.ListenAndServe(servePort, nil))
}

// Open a URL in the default web browser
func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll, FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}
