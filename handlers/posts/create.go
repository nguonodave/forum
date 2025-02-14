package posts

import "net/http"


func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// 20 MB limit
		maxSizeErr := r.ParseMultipartForm(20 << 20)
		if maxSizeErr != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}
	}
}