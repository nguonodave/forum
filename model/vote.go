package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"forum/database"
	"forum/xerrors"

	"github.com/google/uuid"
)

// VoteRequest represents the incoming vote request
type VoteRequest struct {
	UserID    uuid.UUID  `json:"userId"`
	PostID    *uuid.UUID `json:"postId,omitempty"`
	CommentID *uuid.UUID `json:"commentId,omitempty"`
	Type      string     `json:"type"`
}

// VoteResponse represents the response after a vote operation
type VoteResponse struct {
	Success      bool    `json:"success"`
	LikeCount    int     `json:"likeCount"`
	DislikeCount int     `json:"dislikeCount"`
	UserVote     *string `json:"userVote,omitempty"` // "like", "dislike", or null
	Error        string  `json:"error,omitempty"`
}

// HandleVote processes a vote action
func HandleVote(req *VoteRequest) (*VoteResponse, error) {
	// Validate request
	if req.UserID == uuid.Nil {
		return nil, xerrors.ErrInvalidUser
	}
	if req.PostID == nil && req.CommentID == nil {
		return nil, xerrors.ErrInvalidRequest
	}
	if req.Type != "like" && req.Type != "dislike" {
		return nil, xerrors.ErrInvalidVoteType
	}

	fmt.Printf("%+v,\n", req)
	// Start transaction
	tx, err := database.Db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Check for existing vote
	var existingVoteID uuid.UUID
	var existingVoteType string
	var query string
	var args []interface{}

	if req.PostID != nil {
		query = `SELECT id, type FROM votes WHERE user_id = ? AND post_id = ? AND comment_id IS NULL`
		args = []interface{}{req.UserID, *req.PostID}
	} else {
		query = `SELECT id, type FROM votes WHERE user_id = ? AND comment_id = ? AND post_id IS NULL`
		args = []interface{}{req.UserID, *req.CommentID}
	}

	err = tx.QueryRow(query, args...).Scan(&existingVoteID, &existingVoteType)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check existing vote: %v", err)
	}

	// Handle vote based on existing state
	if errors.Is(err, sql.ErrNoRows) {
		// Create new vote
		newVote := uuid.New().String()
		query = `
			INSERT INTO votes (id, user_id, post_id, comment_id, type, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`
		_, err = tx.Exec(query, newVote, req.UserID, req.PostID, req.CommentID, req.Type, time.Now().Format(time.ANSIC))
		if err != nil {
			return nil, fmt.Errorf("failed to create vote: %v", err)
		}
	} else if existingVoteType == req.Type {
		// Remove vote if same type
		query = `DELETE FROM votes WHERE id = ?`
		_, err = tx.Exec(query, existingVoteID)
		if err != nil {
			return nil, fmt.Errorf("failed to remove vote: %v", err)
		}
	} else {
		// Update vote type if different
		query = `UPDATE votes SET type = ? WHERE id = ?`
		_, err = tx.Exec(query, req.Type, existingVoteID)
		if err != nil {
			return nil, fmt.Errorf("failed to update vote: %v", err)
		}
	}

	// Get updated vote counts
	var likeCount, dislikeCount sql.NullInt64
	var userVote sql.NullString

	if req.PostID != nil {
		query = `
			SELECT 
				SUM(CASE WHEN type = 'like' THEN 1 ELSE 0 END) as like_count,
				SUM(CASE WHEN type = 'dislike' THEN 1 ELSE 0 END) as dislike_count,
				(SELECT type FROM votes WHERE user_id = ? AND post_id = ? AND comment_id IS NULL) as user_vote
			FROM votes 
			WHERE post_id = ? AND comment_id IS NULL
		`
		err = tx.QueryRow(query, req.UserID, *req.PostID, *req.PostID).Scan(&likeCount, &dislikeCount, &userVote)
	} else {
		query = `
			SELECT 
				SUM(CASE WHEN type = 'like' THEN 1 ELSE 0 END) as like_count,
				SUM(CASE WHEN type = 'dislike' THEN 1 ELSE 0 END) as dislike_count,
				(SELECT type FROM votes WHERE user_id = ? AND comment_id = ? AND post_id IS NULL) as user_vote
			FROM votes 
			WHERE comment_id = ? AND post_id IS NULL
		`
		err = tx.QueryRow(query, req.UserID, *req.CommentID, *req.CommentID).Scan(&likeCount, &dislikeCount, &userVote)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get vote counts: %v", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Prepare response
	response := &VoteResponse{
		Success:      true,
		LikeCount:    int(likeCount.Int64),
		DislikeCount: int(dislikeCount.Int64),
	}

	if userVote.Valid {
		response.UserVote = &userVote.String
	}

	return response, nil
}
