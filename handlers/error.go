package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

// ErrorPage renders an HTML error page with the provided error message and HTTP status code.
//
// Parameters:
//   - w: The response writer to send the rendered HTML error page
//   - errorText: The error message to display on the page
//   - statusCode: The HTTP status code to set in the response
//
// The function sets the HTTP status code, then renders an error page
// with the provided error message and status code. If there are any errors during template
// parsing or execution, it returns a generic 500 Internal Server Error response.
//
// Example usage:
//
//	ErrorPage(w, "Resource not found", http.StatusNotFound)
func ErrorPage(w http.ResponseWriter, errorText string, statusCode int) {
	type ErrorContent struct {
		Message string
		Code    string
	}

	w.WriteHeader(statusCode)
	content := ErrorContent{
		Message: errorText,
		Code:    strconv.Itoa(statusCode),
	}

	TemplateError := func(message string, err error) {
		http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
		log.Printf("%s: %v", message, err)
	}

	temp, err := template.ParseFiles(filepath.Join(templatesDir, "error.html"))
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
