package utils

import (
	"sort"
	"time"
	"real_time_forum/backend/db/models"
)

func SortPostsByLatest(postsMap map[string]*models.Post) []models.Post {
	// Convert map to slice
	var postsSlice []models.Post
	for _, post := range postsMap {
		post.CommentsNum = len(post.Comments)
		postsSlice = append(postsSlice, *post)
	}

	// Sort by createdAt in descending order
	sort.Slice(postsSlice, func(i, j int) bool {
		// Convert createdAt strings to time.Time
		timeI, _ := time.Parse("2006-01-02 15:04:05", postsSlice[i].CreatedAt)
		timeJ, _ := time.Parse("2006-01-02 15:04:05", postsSlice[j].CreatedAt)
		// Compare times (latest first)
		return timeI.After(timeJ)
	})

	return postsSlice
}
