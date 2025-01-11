package models

import (
	"database/sql"
	"log"
)

type PostRequest struct {
	UserID     int
	Title      string `json:"title"`
	Content    string `json:"content"`
	Categories []int  `json:"categories"`
}

type Post struct {
	ID            int        `json:"id"`
	UserID        int        `json:"user_id"`
	UserFirstName string     `json:"user_first_name"`
	UserLastName  string     `json:"user_last_name"`
	UserNickname  string     `json:"user_nickname"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	CreatedAt     string     `json:"created_at"`
	LikesCount    int        `json:"likes_count"`
	DislikesCount int        `json:"dislike_count"`
	CommentsCount int        `json:"comments_count"`
	Categories    []Category `json:"categories"`
}

type Reaction struct {
	UserID int
	PostID int    `json:"post_id"`
	Type   string `json:"reaction"`
}

func FetchPosts(limit, page int) ([]Post, error) {
	var (
		posts  []Post
		offset = page * limit
	)

	// 'COALESCE' is used to replace NULL values with 0 in a concise way.
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
			COALESCE(like_count, 0) AS likes_count,
			COALESCE(dislike_count, 0) AS dislikes_count,
			COALESCE(comments_count, 0) AS comments_count
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN (
					SELECT post_id, 
						SUM(reaction = 'like') AS like_count,
						SUM(reaction = 'dislike') AS dislike_count
					FROM post_reactions
					GROUP BY post_id
				) reactions ON reactions.post_id = p.id
		LEFT JOIN (
					SELECT post_id, COUNT(id) AS comments_count
					FROM comments
					GROUP BY post_id
				) comments ON comments.post_id = p.id
		ORDER BY p.created_at DESC
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
			&post.LikesCount,
			&post.DislikesCount,
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

func CheckPostExist(postID int) error {
	var id int
	err := DB.QueryRow("SELECT id FROM posts WHERE id = ?", postID).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func StorePost(post PostRequest) (int64, error) {
	query := `INSERT INTO posts (user_id, title, content) VALUES (?,?,?)`

	result, err := DB.Exec(query, post.UserID, post.Title, post.Content)
	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return postID, nil
}

func StorePostCategory(postID int64, categoryID int) (int64, error) {
	query := `INSERT INTO posts_categories (post_id, category_id) VALUES (?,?)`

	result, err := DB.Exec(query, postID, categoryID)
	if err != nil {
		return 0, err
	}

	postcatID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return postcatID, nil
}

func ReactToPost(reaction Reaction) error {
	var oldReaction string

	query := `SELECT reaction FROM post_reactions WHERE post_id = ? AND user_id = ?;`
	err := DB.QueryRow(query, reaction.PostID, reaction.UserID).Scan(&oldReaction)

	if err != nil {
		if err == sql.ErrNoRows {
			query := `INSERT INTO post_reactions (post_id, user_id, reaction) VALUES (?, ?, ?)`
			_, err = DB.Exec(query, reaction.PostID, reaction.UserID, reaction.Type)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if oldReaction == reaction.Type {
			query := `DELETE FROM post_reactions WHERE post_id = ? AND user_id = ?`
			_, err = DB.Exec(query, reaction.PostID, reaction.UserID)
			if err != nil {
				return err
			}
		} else {
			query := `UPDATE post_reactions SET reaction = ? WHERE post_id = ? AND user_id = ?`
			_, err = DB.Exec(query, reaction.Type, reaction.PostID, reaction.UserID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
