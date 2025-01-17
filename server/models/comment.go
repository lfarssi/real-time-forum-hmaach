package models

type CommentRequest struct {
	UserID  int    `json:"user_id"`
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
}

type Comment struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	PostID        int    `json:"post_id"`
	UserFirstName string `json:"first_name"`
	UserLastName  string `json:"last_name"`
	UserNickname  string `json:"nickname"`
	Content       string `json:"content"`
	CreatedAt     string `json:"created_at"`
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

func StoreComment(comment CommentRequest) (int64, error) {
	query := `INSERT INTO comments (user_id, post_id, content) VALUES (?,?,?)`

	result, err := DB.Exec(query, comment.UserID, comment.PostID, comment.Content)
	if err != nil {
		return 0, err
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return commentID, nil
}
