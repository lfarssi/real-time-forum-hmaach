package models

import (
	"fmt"
	"log"
)

type Post struct {
	ID            int
	UserID        int
	UserFirstName string
	UserLastName  string
	UserNickname  string
	Title         string
	Content       string
	CreatedAt     string
	CommentsCount int
	Categories    []Category
}

func FetchPosts(limit, page int) ([]Post, error) {
	var (
		posts  []Post
		offset = page * limit
	)

	query := `
	SELECT
		p.id,
		p.user_id,
		u.first_name, 
		u.last_name, 
		u.nickname, 
		p.title,
		p.content,
		p.created_at,
		(
			SELECT
				COUNT(c.id)
			FROM
				comments c
			WHERE
				c.post_id = p.id
		) AS comments_count
	FROM
		posts p
	INNER JOIN users u ON p.user_id = u.id
	ORDER BY
		p.created_at DESC
	LIMIT ? OFFSET ?;`

	rows, err := DB.Query(query, limit, offset)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.UserFirstName,
			&post.UserLastName,
			&post.UserNickname,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.CommentsCount,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		// Fetch categories for the post
		post.Categories, err = FetchCategoriesByPostID(post.ID)
		if err != nil {
			log.Println("Error fetching categories for post:", post.ID, err)
			return nil, err
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func StorePost(user_id int, title, content string) (int64, error) {
	query := `INSERT INTO posts (user_id,title,content) VALUES (?,?,?)`

	result, err := DB.Exec(query, user_id, title, content)
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}

	postID, _ := result.LastInsertId()

	return postID, nil
}

func StorePostCategory(post_id int64, category_id int) (int64, error) {
	query := `INSERT INTO post_category (post_id, category_id) VALUES (?,?)`

	result, err := DB.Exec(query, post_id, category_id)
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}

	postcatID, _ := result.LastInsertId()

	return postcatID, nil
}
