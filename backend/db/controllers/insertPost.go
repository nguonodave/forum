package controllers

import (
	// "errors"

	"errors"

	"github.com/gofrs/uuid"
)


func AddNewPostToDB(id uuid.UUID, userId string, userName string, title string, content string, createdAt string, category string) error {
	if title == "" || content == "" {
		return errors.New("All fields must not be empty")
	}
	stmt := `INSERT INTO posts(id, userid, username, title, content, createdAt, category, likes, dislikes)
	VALUES (?,?,?,?,?,?,?,?,?)`
	_, err := DB.Exec(stmt, id.String(), userId, userName, title,  content, createdAt, category, 0, 0)	
	if err != nil {
		return err
	}
	return nil
}