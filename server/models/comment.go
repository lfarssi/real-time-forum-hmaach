package models

import (
	"fmt"
)

type Comment struct {
	ID        int
	UserID    int
	PostID    int
	nickname  string
	Content   string
	Likes     int
	Dislikes  int
	CreatedAt string
}

func FetchCommentsByPostID(postID int) ([]Comment, error) {
	var comments []Comment
	query := `
	SELECT
		c.id,
		c.user_id,
		u.nickname,
		c.content,
		strftime('%m/%d/%Y %I:%M %p', c.created_at) AS formatted_created_at,
	FROM
		comments c
	INNER JOIN users u 
	ON c.user_id = u.id
	WHERE
		c.post_id = ?
	ORDER BY
		c.created_at DESC
	`

	rows, err := DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.nickname,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		comment.PostID = postID

		comments = append(comments, comment)
	}

	return comments, nil
}

func StoreComment(user_id, post_id int, content string) (int64, error) {
	query := `INSERT INTO comments (user_id, post_id, content) VALUES (?,?,?)`

	result, err := DB.Exec(query, user_id, post_id, content)
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}

	commentID, _ := result.LastInsertId()

	return commentID, nil
}

// Count comments by post ID
func CountCommentsByPostID(postID int) (int, error) {
	var count int
	query := "SELECT COUNT(id) FROM comments WHERE post_id = ?"
	err := DB.QueryRow(query, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting comments: %v", err)
	}
	return count, nil
}
