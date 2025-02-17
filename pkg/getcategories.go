package pkg

import (
	"log"
	"net/http"

	"forum/database"
)

type category struct {
	Id   string
	Name string
}

func GetCategories(w http.ResponseWriter) ([]category, error) {
	// fetch all categories to render to the create post form
	categRows, categQueryErr := database.Db.Query(`SELECT id, name FROM categories`)
	if categQueryErr != nil {
		log.Printf("Error fetching categories: %v\n", categQueryErr)
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return nil, categQueryErr
	}
	defer categRows.Close()

	var categories []category

	for categRows.Next() {
		var category struct {
			Id   string
			Name string
		}
		err := categRows.Scan(&category.Id, &category.Name)
		if err != nil {
			log.Printf("Error scanning category: %v\n", err)
			continue
		}
		categories = append(categories, category)
	}

	return categories, nil
}
