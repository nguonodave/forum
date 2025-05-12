package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"real_time_forum/backend/db/controllers"
	"real_time_forum/backend/db/models"
	"strings"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		request, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		user := models.User{}

		if err := json.Unmarshal(request, &user); err != nil {
			log.Println(err)
			http.Error(w, "Invalid request payload", http.StatusInternalServerError)
			return
		}
		
		var username = user.Username
		var email = user.Email
		var password = user.Password
		var age = user.Age
		var firstName = user.FirstName
		var lastName = user.LastName
		var gender = user.Gender

		email = strings.TrimSpace(email)
		firstName = strings.TrimSpace(firstName)
		lastName = strings.TrimSpace(lastName)
		username = strings.TrimSpace(username)
		password = strings.TrimSpace(password)
		gender = strings.TrimSpace(gender)

		if email == "" || firstName == "" || lastName == "" || username == "" || password == "" || age == 0 || gender == "" {
			http.Error(w, "cannot enter empty values", http.StatusBadRequest)
			return
		}

		userID, err := uuid.NewV4()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil {
			http.Error(w, "internal sever error", http.StatusInternalServerError)
			return
		}
		err = controllers.AddNewUsersToDB(userID,  firstName, lastName, username, gender, email, string(passwordHash), age)
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Account Successfully Created"})
	}

}
