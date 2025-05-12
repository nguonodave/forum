package controllers

import (
	"database/sql"
	"log"

	"github.com/gofrs/uuid"
)

func IsThread(user1Id, user2Id string) (string, bool, error) {
	// log.Printf("Checking thread between %s and %s", user1Id, user2Id)

	stmt := `SELECT threadId FROM messages WHERE 
	         (user1Id = ? AND user2Id = ?) OR 
	         (user1Id = ? AND user2Id = ?) LIMIT 1`

	row := DB.QueryRow(stmt, user1Id, user2Id, user2Id, user1Id)

	var threadId string
	err := row.Scan(&threadId)
	// log.Println("TREAD ID IN IS THREAD FUNC IS: ", threadId)
	if err != nil {
		if err == sql.ErrNoRows {
			// log.Println("No thread found")
			return "", false, nil
		}
		log.Println("Error querying thread:", err)
		return "", false, err
	}

	// log.Printf("Found thread with ID: %s", threadId)
	return threadId, true, nil
}

func CreateNewMessageTable(user1Id, user2Id string) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	threadId, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	stm := `INSERT INTO messages(id, user1Id, user2Id, threadId) VALUES(?,?,?,?)`

	_, err = DB.Exec(stm, id.String(), user1Id, user2Id, threadId.String(),)

	if err != nil {
		return "", err
	}

	return threadId.String(), nil
}

func CreateNewThread(message, id string) error {
	stm := `INSERT INTO thread(id, messages) VALUES(?, ?) `
	_, err = DB.Exec(stm, id, message)

	if err != nil {
		return err
	}

	return nil
}


func GetMessages(threadId string) (string, error) {
	var messages string

	stm := `SELECT messages FROM thread WHERE id = ?`

	err := DB.QueryRow(stm, threadId).Scan(&messages)

	if err != nil {
		return "", err
	}

	return messages, nil
}

func UpdateMessageThread(threadId, message string) error {
	stm := "UPDATE thread SET messages = ? WHERE id = ?"

	_, err := DB.Exec(stm, message, threadId)
	if err != nil {
		return err
	}

	return nil
}

func GetUserThreadIds(userid string) ([]string, error) {
	var (
		threadIds []string
		threadId  string
	)

	// Corrected SQL statement
	stm := `SELECT threadId FROM messages WHERE (user1Id = ? OR user2Id = ?)`

	rows, err := DB.Query(stm, userid, userid)
	if err != nil {
		return threadIds, err
	}
	defer rows.Close() // Ensure rows are closed after iteration

	for rows.Next() {
		err = rows.Scan(&threadId)
		if err != nil {
			log.Println(err)
			return threadIds, err
		}
		threadIds = append(threadIds, threadId)
	}

	// Check for errors after iteration
	if err = rows.Err(); err != nil {
		return threadIds, err
	}

	return threadIds, nil
}
