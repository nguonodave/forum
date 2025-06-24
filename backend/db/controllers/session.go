package controllers

import (
	"database/sql"
	"log"

	"github.com/gofrs/uuid"
)

// StoreSession stores a new session in the database
func StoreSession(sessionToken string, userID string, time string) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	stm, err := DB.Prepare(`INSERT INTO sessions (id, sessionToken, userid, expires)
	VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(id.String(), sessionToken, userID, time)
	if err != nil {
		return err
	}

	return nil
}

// ValidSession checks if a session token exists in the database
func ValidSession(cookie string) bool {
	stm, err := DB.Prepare("SELECT id FROM sessions WHERE sessionToken = ?")
	if err != nil {
		log.Println("Failed to prepare query:", err)
		return false
	}
	defer stm.Close()

	var id string
	err = stm.QueryRow(cookie).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false 
		}
		log.Println("Query error:", err)
		return false
	}

	return true
}

// GetUserDetailsFromSession retrieves user details associated with a session
func GetUserDetailsFromSession(cookie string) (string, string, string, error) {
	var userId, username, email, sessionId string

	stm := `
SELECT
    sessions.userid AS sessionUserId,
    users.id AS userId,
    users.username AS username,
    users.email AS userEmail
FROM
    sessions
JOIN
    users
ON
    users.id = sessions.userid
WHERE
    sessions.sessionToken = ?
`

	row := DB.QueryRow(stm, cookie)
	err := row.Scan(&sessionId, &userId, &username, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", "", err // No matching session found
		}
		return "", "", "", err
	}

	return userId, username, email, nil
}

func EndSession(cookie string) error {

	stm, err := DB.Prepare("DELETE FROM sessions WHERE sessionToken = ?")
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(cookie)
	if err != nil {
		return err
	}

	return nil
}

func CheckExistingSessions(userid string) bool {
	var session_token sql.NullString
	stm, err := DB.Prepare("SELECT sessionToken FROM sessions WHERE userid = ?")
	if err != nil {
		return false
	}

	defer stm.Close()

	err = stm.QueryRow(userid).Scan(&session_token)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		return false
	}

	return true
}

func DeleteExistingSession(userid string) error {

	stm, err := DB.Prepare("DELETE FROM sessions WHERE userid = ?")
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(userid)
	if err != nil {
		return err
	}

	return nil
}

