package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Index handler designed for the application's index page
func Index(w http.ResponseWriter, r *http.Request) {
	templateFile := "base.html"
	type Content struct {
		Message string
		Data    string
	}

	content := Content{
		Message: "Some message to pass to template",
		Data:    "Some data to pass to template",
	}

	TemplateError := func(message string, err error) {
		http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
		log.Printf("%s: %v", message, err)
	}
	temp, err := template.ParseFiles(filepath.Join(templatesDir, templateFile))
	if err != nil {
		TemplateError("error parsing template", err)
		return
	}

	err = temp.Execute(w, content)
	if err != nil {
		TemplateError("error executing template", err)
		return
	}

}
