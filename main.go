package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	"forum/controller"

	"forum/middlewares"

	"forum/database"
	"forum/fileio"
	"forum/handlers"
)

var (
	port = flag.Int("P", 8080, "port to listen on")
	open = flag.Bool("O", false, "open server index page in the default browser")
)

func main() {
	flag.Parse()

	{
		logFilePath := path.Join(os.TempDir(), fmt.Sprintf("%d-forum-logger.log", os.Getpid()))
		logger, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			log.Printf("failed to setup file logging: logging to stderr instead: %v\n", err)
		} else {
			log.Printf("saving logs to: %s\n", logFilePath)
		}
		w := io.MultiWriter(os.Stdout, logger)
		log.SetOutput(w)
		defer fileio.Close(logger)
	}

	err := database.InitializeDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.Db.Close()
	fmt.Println("Database initialized successfully!")

	var parseTemplateErr error
	handlers.Templates, parseTemplateErr = handlers.ParseTemplates()
	if parseTemplateErr != nil {
		log.Fatalf("Failed to parse templates: %v", parseTemplateErr)
	}

	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/login", middlewares.RedirectIfLoggedIn(handlers.Login))
	http.HandleFunc("/register", middlewares.RedirectIfLoggedIn(handlers.Register))
	http.HandleFunc("/api/vote", controller.ValidateSession(handlers.HandleVoteRequest))
	http.HandleFunc("/logout", controller.ValidateSession(handlers.Logout))
	http.HandleFunc("/api/comment", controller.ValidateSession(handlers.AddCommentHandler(database.Db)))

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

	servePort := fmt.Sprintf(":%d", *port)
	url := fmt.Sprintf("http://localhost%s", servePort)
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
