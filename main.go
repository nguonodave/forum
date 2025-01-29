package main

import (
	"fmt"
	"log"

	"forum/database"
)

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

