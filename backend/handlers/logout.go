package handlers


import (
	"real_time_forum/backend/db/controllers"
	"net/http"
	"time"
	"encoding/json"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the session token from the cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "No session token provided"}) 
		return
	}

	// End the session in the database
	err = controllers.EndSession(cookie.Value)
	if err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()}) 
		return
	}

	// Clear the session cookie on the client side
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // Set expiration to the past
		HttpOnly: true,
	})
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "user logged out sucessfully"}) 
}
