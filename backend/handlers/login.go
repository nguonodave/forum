package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"time"

	"real_time_forum/backend/db/controllers"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Credential struct {
	Name string `json:"name"`
	Password string `json:"password"`
}


func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var creds Credential

	request, err := io.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(request, &creds)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	email := creds.Name
	password := creds.Password
	
	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}
	
	// Verify user credentials
	storedPassword, userID, useremail, username, firstname, lastname, gender, age, err := controllers.VerifyUser(email, password)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message" : "User Not Found. Please Create Account To Continue"})
		return		
	}
	if !checkPasswordHash(password, storedPassword) {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message" : "Inavalid UserName/password"})
		return
	}
	
	if controllers.CheckExistingSessions(userID) {
		controllers.DeleteExistingSession(userID)
	}

	// log.Println(email, userID, password)
	sessionToken := generateToken(32)

	err = controllers.StoreSession(sessionToken, userID, time.Now().Add(24*time.Hour).String())

	if err != nil {
		log.Println("Failed to store session:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path: "/",
	})


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message" : "logged In",
		"session_token" : sessionToken,
		"id" : userID, 
		"email" : useremail, 
		"username" : username, 
		"firstname" : firstname, 
		"lastname" : lastname, 
		"gender" : gender, 
		"age" : age,
	})
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
