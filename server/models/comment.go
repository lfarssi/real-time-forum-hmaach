package models

import (
	"fmt"
)

type Comment struct {
	ID            int
	UserID        int
	PostID        int
	UserFirstName string
	UserLastName  string
	UserNickname  string
	Content       string
	CreatedAt     string
}

func FetchCommentsByPostID(postID, limit, page int) ([]Comment, error) {
	var (
		comments []Comment
		offset   = page * limit
	)

	query := `
	SELECT
		c.id,
		c.user_id,
		u.first_name, 
		u.last_name, 
		u.nickname, 
		c.content,
		c.created_at
	FROM
		comments c
	INNER JOIN users u 	ON c.user_id = u.id
	WHERE
		c.post_id = ?
	ORDER BY
		c.created_at DESC
	LIMIT ? OFFSET ?;
	`

	rows, err := DB.Query(query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.UserFirstName,
			&comment.UserLastName,
			&comment.UserNickname,
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
