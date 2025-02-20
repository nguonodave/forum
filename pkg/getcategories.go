package pkg

import (
	"fmt"
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
	categRows, categQueryErr := database.Db.Query(`SELECT name FROM categories`)
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
		err := categRows.Scan(&category.Name)
		if err != nil {
			log.Printf("Error scanning category: %v\n", err)
			continue
		}
		categories = append(categories, category)
	}
	fmt.Println("categoreis", categories)
	return categories, nil
}

func ValidCategory(category string, categories []category) bool {
	for _, v := range categories {
		if category == v.Name {
			return true
		}
	}
	return false
}
