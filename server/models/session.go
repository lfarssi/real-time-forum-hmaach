package models

import (
	"database/sql"
	"log"
	"strings"
	"time"
)

func StoreSession(userID int, tokenID string, expires_at time.Time) error {
	query := `INSERT OR REPLACE INTO sessions (user_id, token, expires_at) VALUES (?,?,?)`

	_, err := DB.Exec(query, userID, tokenID, expires_at)
	if err != nil {
		return err
	}

	return nil
}

// ValidSession validates the session by checking the token in the database and verifying it has not expired.
func ValidSession(token string) (int, bool, string) {
	var (
		expiration time.Time
		userID     int
	)
	token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))

	// Query the sessions table to get the user and expiration time
	query := `
		SELECT user_id, expires_at 
		FROM sessions 
		WHERE token = ?
	`
	row := DB.QueryRow(query, token)

	// Scan the result into variables
	err := row.Scan(&userID, &expiration)
	if err == sql.ErrNoRows {
		return 0, false, "unauthorized"
	} else if err != nil {
		log.Println("Database error:", err)
		return 0, false, "Internal Server Error"
	}

	// Check if the session is expired
	if time.Now().After(expiration) {
		return 0, false, "unauthorized"
	}

	// Session is valid
	return userID, true, "success"
}

func DeleteUserSession(userID int) error {
	_, err := DB.Exec(`DELETE FROM sessions WHERE user_id = ?;`, userID)
	return err
}
