package controllers

import (
"github.com/gofrs/uuid"
"fmt"
sqlite3"github.com/mattn/go-sqlite3"
)


func AddNewUsersToDB(id uuid.UUID,  firstName, lastName, nickname, gender, email, password string, age int) error {
	stmt := `INSERT INTO users(id, firstname, lastname, username, gender, age, email, password)
	VALUES (?,?,?,?,?,?,?,?)`
	_, err := DB.Exec(stmt, id.String(), firstName, lastName, nickname, gender, age, email, password)	
	if err != nil {
		// Check for unique constraint violation (SQLite error)
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("the email address '%s' is already taken", email)
		}
		return err
	}
	return nil
}

func InsertCommentToDB(id string, postId string, userId string,username string, context string, createdAt string) (int, error) {
	stmt := `INSERT INTO comments(id, userid, username , postid, content, createdAt, likes, dislikes)
	VALUES (?,?,?,?,?,?,?,?)`
	_, err := DB.Exec(stmt, id, userId, username, postId, context, createdAt, 0, 0)	
	if err != nil {
		return 0, err
	}

	comment_count, err := CommentCount(postId)
	if err != nil {
		return 0, err
	}
	return comment_count, nil
}

func CommentCount(postId string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM comments WHERE postid = ?`
	err := DB.QueryRow(query, postId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}