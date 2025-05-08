package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func VerifyUser(email, password string) (string, string, string, string, string, string, string, string, error) {
	var userID, storedPassword, useremail, username, firstname, lastname, gender, age string

	var stats string
	if strings.Contains(email, "@") {
		stats = "email"
	} else {
		stats = "username"
	}

	stm := fmt.Sprintf("SELECT id, username, firstname, lastname, gender, age, email, password FROM users WHERE %s = ?",stats)

	err = DB.QueryRow(stm, email).Scan(&userID, &username, &firstname, &lastname, &gender, &age, &useremail, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", "", "","", "", "", "", err
		} else {
			log.Fatal(err)
		}
	}
	return storedPassword, userID, useremail, username, firstname, lastname, gender, age, nil
}
