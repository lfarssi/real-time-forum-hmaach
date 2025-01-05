package models

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type Post struct {
	ID            int
	UserID        int
	UserNickname  string
	Title         string
	Content       string
	CreatedAt     string
	Comments      int
	CategoriesStr string
	Categories    []string
}

type PostDetail struct {
	Post     Post
	Comments []Comment
}

func FetchPosts(page int) ([]Post, error) {
	var posts []Post
	// Query to fetch posts
	// nickname
	query := `SELECT
		p.id,
		p.user_id,
		u.nickname, 
		p.title,
		p.content,
		strftime('%m/%d/%Y %I:%M %p', p.created_at) AS formatted_created_at,
		(
			SELECT
				COUNT(*)
			FROM
				comments c
			WHERE
				c.post_id = p.id
		) AS comments_count,
		(
			SELECT
				GROUP_CONCAT(c.label)
			FROM
				categories c
			INNER JOIN post_category pc ON c.id = pc.category_id
			WHERE
				pc.post_id = p.id
		) AS categories
	FROM
		posts p
		INNER JOIN users u ON p.user_id = u.id
	ORDER BY
		p.created_at DESC
	LIMIT 10 OFFSET ? ;
	`
	rows, err := DB.Query(query, page)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows
	for rows.Next() {
		var post Post
		// Scan the data into the Post struct
		err := rows.Scan(&post.ID,
			&post.UserID,
			&post.UserNickname,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Comments,
			&post.CategoriesStr)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		// it came from the  database as "technology,sports...", so we need to split it
		post.Categories = strings.Split(post.CategoriesStr, ",")

		// Append the Post struct to the posts slice
		posts = append(posts, post)
	}
	return posts, nil
}

func FetchPost(postID int) (PostDetail, error) {
	var post Post
	post.ID = postID

	// Query to fetch the post
	query := `SELECT
		p.user_id,
		u.nickname,
		p.title,
		p.content,
		strftime('%m/%d/%Y %I:%M %p', p.created_at) AS formatted_created_at,
		(
			SELECT COUNT(*)
			FROM comments c
			WHERE c.post_id = p.id
		) AS comments_count,
		(
			SELECT GROUP_CONCAT(c.label)
			FROM categories c
			INNER JOIN post_category pc ON c.id = pc.category_id
			WHERE pc.post_id = p.id
		) AS categories
	FROM
		posts p
		INNER JOIN users u ON p.user_id = u.id
	WHERE p.id = ?`

	// Use QueryRow for a single result
	row := DB.QueryRow(query, postID)

	// Scan the data into the Post struct
	err := row.Scan(
		&post.UserID,
		&post.UserNickname,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.Comments,
		&post.CategoriesStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return PostDetail{}, err
		}
		log.Println("Error scanning row:", err)
		return PostDetail{}, err
	}

	// Process categories
	post.Categories = strings.Split(post.CategoriesStr, ",")

	// Format the created_at field
	// post.CreatedAt = post.CreatedAt.Format("01/02/2006 03:04 PM")
	comments, err := FetchCommentsByPostID(postID)
	if err != nil {
		log.Println("Error fetching comments from the database:", err)
	}

	return PostDetail{
		Post:     post,
		Comments: comments,
	}, nil
}

func FetchPostsByCategory(categoryID, page int) ([]Post, error) {
	var posts []Post
	query := `
		SELECT
			p.id,
			p.user_id,
			u.nickname,
			p.title,
			p.content,
			strftime('%m/%d/%Y %I:%M %p', p.created_at) AS formatted_created_at,
			(
				SELECT
					COUNT(*)
				FROM
					comments c
				WHERE
					c.post_id = p.id
			) AS comments_count,
			(
				SELECT
					GROUP_CONCAT(c.label)
				FROM
					categories c
				INNER JOIN post_category pc ON c.id = pc.category_id
				WHERE
					pc.post_id = p.id
			) AS categories
		FROM
			posts p
			INNER JOIN users u ON p.user_id = u.id
			INNER JOIN post_category pc ON p.id = pc.post_id
		WHERE pc.category_id = ?
		ORDER BY
			p.created_at
		LIMIT 10 OFFSET ? ;
	`
	rows, err := DB.Query(query, categoryID, page)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID,
			&post.UserID,
			&post.UserNickname,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Comments,
			&post.CategoriesStr)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		// it came from the  database as "technology,sports...", so we need to split it
		post.Categories = strings.Split(post.CategoriesStr, ",")

		posts = append(posts, post)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
		return nil, err
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

