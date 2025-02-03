package main

import (
	"flag"
	"fmt"
	"forum/database"
	"forum/fileio"
	"forum/handlers"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
)

var port = flag.Int("P", 8080, "port to listen on")
var open = flag.Bool("O", false, "open server index page in the default browser")

func main() {
	println("hello forum")
	// Initialize database
	db, err := database.InitializeDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	fmt.Println("Database operations completed successfully!")
}

// openBrowser opens a URL in the default web browser based on the operating
// system that the code is running on. It handles Linux, Windows,and macOS platforms.
// It takes a single parameter which is a string representing the URL to open.
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
		log.Printf("Failed to open browser:%v", err)
	}
}
