package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
)

func CheckPreviousReaction(postid, userid string) (string, error) {
	stm , err := DB.Prepare("SELECT reaction FROM likes WHERE postid = ? AND userid = ?")
	if err != nil {
		return "", err
	}
	defer stm.Close()
	
	var react string
	err = stm.QueryRow(postid, userid).Scan(&react)
    if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		} 
        return "", err // Return the error if the query execution fails
    }

	return react, nil
}

func UpdateReaction(newReaction string, postid string, userid string) error {
	stm, err := DB.Prepare("UPDATE likes SET reaction = ? WHERE postid = ? AND userid = ?")
    if err != nil {
        return err 
    }
    defer stm.Close() 

    _, err = stm.Exec(newReaction, postid, userid)
    if err != nil {
        return err 
    }

    return nil 

}


func DeleteReaction(postId, userId string) error {
	// Prepare the DELETE query
	query := `DELETE FROM likes WHERE postid = ? AND userid = ?`

	// Execute the query
	result, err := DB.Exec(query, postId, userId)
	if err != nil {
		return fmt.Errorf("failed to delete row: %v", err)
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows found with postId %s and userId %s", postId, userId)
	}

	// fmt.Printf("Successfully deleted %d row(s)\n", rowsAffected)
	return nil
}

func DecreaseLikes(postId string) (int, error) {
	// Step 1: Get the current number of likes
	var likes int
	query := `SELECT likes FROM posts WHERE id = ?`
	err := DB.QueryRow(query, postId).Scan(&likes)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch likes for post ID %s: %v", postId, err)
	}

	// Step 2: Update the likes count
	updateQuery := `UPDATE posts SET likes = ? WHERE id = ?`
	result, err := DB.Exec(updateQuery, likes-1, postId)
	if err != nil {
		return 0, fmt.Errorf("failed to update likes for post ID %s: %v", postId, err)
	}

	// Step 3: Validate the update operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return  0, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return 0, fmt.Errorf("no rows updated for post ID %s", postId)
	}

	// fmt.Printf("Successfully decreased likes for post ID %s. New likes count: %d\n", postId, likes-1)
	return likes-1, nil
}


func DecreaseDislikes(postId string) (int, error) {
	// Step 1: Get the current number of likes
	var likes int
	query := `SELECT dislikes FROM posts WHERE id = ?`
	err := DB.QueryRow(query, postId).Scan(&likes)
	if err != nil {
		return 0,  fmt.Errorf("failed to fetch likes for post ID %s: %v", postId, err)
	}

	// Step 2: Update the likes count
	updateQuery := `UPDATE posts SET dislikes = ? WHERE id = ?`
	result, err := DB.Exec(updateQuery, likes-1, postId)
	if err != nil {
		return 0,  fmt.Errorf("failed to update likes for post ID %s: %v", postId, err)
	}

	// Step 3: Validate the update operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0,  fmt.Errorf("failed to retrieve rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return 0, fmt.Errorf("no rows updated for post ID %s", postId)
	}

	// fmt.Printf("Successfully decreased likes for post ID %s. New likes count: %d\n", postId, likes-1)
	return likes-1, nil
}



func InsertReaction(userId string, postid string, reaction string) error {
	stmt := `INSERT INTO likes(id, userid, postid, reaction)
	VALUES (?,?,?,?)`

	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	
	_, err = DB.Exec(stmt, id.String(), userId, postid, reaction)	
	if err != nil {
		return err
	}
	return nil
}


func IncreaseLikes(postId string) (int, error) {
	// Step 1: Get the current number of likes
	var likes int
	query := `SELECT likes FROM posts WHERE id = ?`
	err := DB.QueryRow(query, postId).Scan(&likes)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch likes for post ID %s: %v", postId, err)
	}
	// fmt.Println(likes)
	// Step 2: Update the likes count
	updateQuery := `UPDATE posts SET likes = ? WHERE id = ?`
	result, err := DB.Exec(updateQuery, likes+1, postId)
	if err != nil {
		return 0, fmt.Errorf("failed to update likes for post ID %s: %v", postId, err)
	}

	// Step 3: Validate the update operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0,  fmt.Errorf("failed to retrieve rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return 0,  fmt.Errorf("no rows updated for post ID %s", postId)
	}

	// fmt.Printf("Successfully increased likes for post ID %s. New likes count: %d\n", postId, likes+1)
	return likes+1, nil
}


func IncreaseDislikes(postId string) (int, error) {
	// Step 1: Get the current number of likes
	var likes int
	query := `SELECT dislikes FROM posts WHERE id = ?`
	err := DB.QueryRow(query, postId).Scan(&likes)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch likes for post ID %s: %v", postId, err)
	}

	// fmt.Println("LIKES >>", likes)

	// Step 2: Update the likes count
	updateQuery := `UPDATE posts SET dislikes = ? WHERE id = ?`
	result, err := DB.Exec(updateQuery, likes+1, postId)
	if err != nil {
		return 0, fmt.Errorf("failed to update likes for post ID %s: %v", postId, err)
	}

	// Step 3: Validate the update operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return 0, fmt.Errorf("no rows updated for post ID %s", postId)
	}

	// fmt.Printf("Successfully increased likes for post ID %s. New likes count: %d\n", postId, likes+1)
	return likes+1, nil
}
