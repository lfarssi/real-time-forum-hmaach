package models

import (
	"fmt"
	"log"
	"strings"
)

type Category struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

func FetchCategories() ([]Category, error) {
	var categories []Category
	query := `
		SELECT
			c.id,
			c.label,
		FROM categories c
		ORDER BY posts_count DESC;
	`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var category Category
		rows.Scan(&category.ID, &category.Label)
		categories = append(categories, category)
	}
	return categories, nil
}

func FetchCategoriesByPostID(postID int) ([]Category, error) {
	query := `
	SELECT 
		c.id,
		c.label
	FROM
		categories c
	INNER JOIN posts_categories pc ON pc.category_id = c.id
	WHERE pc.post_id = ?`

	rows, err := DB.Query(query, postID)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.ID,
			&category.Label,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func CheckCategoriesExist(ids []int) error {
	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(`
        SELECT id
        FROM categories
        WHERE id IN (%s);
    `, placeholders)

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		count++
	}
	if count != len(ids) {
		return fmt.Errorf("categories does not exists in db")
	}

	return nil
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
