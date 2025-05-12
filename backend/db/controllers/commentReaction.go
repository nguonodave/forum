package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
)

func CheckPreviousCommentReaction(commentid, userid string) (string, error) {
	stm , err := DB.Prepare("SELECT reaction FROM commentlikes WHERE commentid = ? AND userid = ?")
	if err != nil {
		return "", err
	}
	defer stm.Close()
	
	var react string
	err = stm.QueryRow(commentid, userid).Scan(&react)
    if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		} 
        return "", err // Return the error if the query execution fails
    }

	return react, nil
}

func UpdateCommentReaction(newReaction string, commentid string, userid string) error {
	stm, err := DB.Prepare("UPDATE commentlikes SET reaction = ? WHERE commentid = ? AND userid = ?")
    if err != nil {
        return err 
    }
    defer stm.Close() 

    _, err = stm.Exec(newReaction, commentid, userid)
    if err != nil {
        return err 
    }

    return nil 

}


func DeleteCommentReaction(commentid, userId string) error {
	// Prepare the DELETE query
	query := `DELETE FROM commentlikes WHERE commentid = ? AND userid = ?`

	// Execute the query
	result, err := DB.Exec(query, commentid, userId)
	if err != nil {
		return fmt.Errorf("failed to delete row: %v", err)
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows found with postId %s and userId %s", commentid, userId)
	}

	return nil
}

func DecreaseCommentLikes(commentid string) error {
	// Step 1: Get the current number of likes
	var likes int
	query := `SELECT likes FROM comments WHERE id = ?`
	err := DB.QueryRow(query, commentid).Scan(&likes)
	if err != nil {
		return fmt.Errorf("failed to fetch likes for post ID %s: %v", commentid, err)
	}

	// Step 2: Update the likes count
	updateQuery := `UPDATE comments SET likes = ? WHERE id = ?`
	result, err := DB.Exec(updateQuery, likes-1, commentid)
	if err != nil {
		return fmt.Errorf("failed to update likes for post ID %s: %v", commentid, err)
	}

	// Step 3: Validate the update operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated for post ID %s", commentid)
	}

	return nil
}


func DecreaseCommentDislikes(commentid string) error {
	// Step 1: Get the current number of likes
	var likes int
	query := `SELECT dislikes FROM comments WHERE id = ?`
	err := DB.QueryRow(query, commentid).Scan(&likes)
	if err != nil {
		return fmt.Errorf("failed to fetch likes for post ID %s: %v", commentid, err)
	}

	// Step 2: Update the likes count
	updateQuery := `UPDATE comments SET dislikes = ? WHERE id = ?`
	result, err := DB.Exec(updateQuery, likes-1, commentid)
	if err != nil {
		return fmt.Errorf("failed to update likes for post ID %s: %v", commentid, err)
	}

	// Step 3: Validate the update operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated for post ID %s", commentid)
	}

	return nil
}



func InsertCommentReaction(userId string, commentid string, reaction string) error {
	stmt := `INSERT INTO commentlikes(id, userid, commentid, reaction)
	VALUES (?,?,?,?)`

	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	
	_, err = DB.Exec(stmt, id.String(), userId, commentid, reaction)	
	if err != nil {
		return err
	}
	return nil
}


func IncreaseCommentLikes(commentid string) error {
	// Step 1: Get the current number of likes
	var likes int
	query := `SELECT likes FROM comments WHERE id = ?`
	err := DB.QueryRow(query, commentid).Scan(&likes)
	if err != nil {
		return fmt.Errorf("failed to fetch likes for post ID %s: %v", commentid, err)
	}
	fmt.Println(likes)
	// Step 2: Update the likes count
	updateQuery := `UPDATE comments SET likes = ? WHERE id = ?`
	result, err := DB.Exec(updateQuery, likes+1, commentid)
	if err != nil {
		return fmt.Errorf("failed to update likes for post ID %s: %v", commentid, err)
	}

	// Step 3: Validate the update operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated for post ID %s", commentid)
	}

	return nil
}


func IncreaseCommentDislikes(commentid string) error {
	// Step 1: Get the current number of likes
	var likes int
	query := `SELECT dislikes FROM comments WHERE id = ?`
	err := DB.QueryRow(query, commentid).Scan(&likes)
	if err != nil {
		return fmt.Errorf("failed to fetch likes for post ID %s: %v", commentid, err)
	}

	// fmt.Println("LIKES >>", likes)

	// Step 2: Update the likes count
	updateQuery := `UPDATE comments SET dislikes = ? WHERE id = ?`
	result, err := DB.Exec(updateQuery, likes+1, commentid)
	if err != nil {
		return fmt.Errorf("failed to update likes for post ID %s: %v", commentid, err)
	}

	// Step 3: Validate the update operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated for post ID %s", commentid)
	}

	return nil
}
