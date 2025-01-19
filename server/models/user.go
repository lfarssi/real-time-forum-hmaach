package models

import (
	"time"

	"forum/server/utils"
)

// represents the data for user registration.
type RegistrationRequest struct {
	FirstName string
	LastName  string
	Email     string
	Nickname  string
	Gender    string
	Age       int
	Password  string
}

// represents the data for user login.
type LoginRequest struct {
	Identifier string
	Password   string
}

type User struct {
	ID          int         `json:"id"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	Nickname    string      `json:"nickname"`
	Email       string      `json:"email"`
	Age         int         `json:"age"`
	Gender      string      `json:"gender"`
	LastMessage LastMessage `json:"last_message,omitempty"`
}
type LastMessage struct {
	Content  string `json:"message"`
	SenderID string `json:"sender_id"`
	SentAt   string `json:"sent_at"`
}

func GenerateSession(userId int) (User, string, error) {
	var user User
	token, err := utils.GenerateToken()
	if err != nil {
		return User{}, "", err
	}

	err = StoreSession(userId, token, time.Now().Add(10*time.Hour))
	if err != nil {
		return User{}, "", err
	}

	user, err = GetUserInfo(userId)
	if err != nil {
		return User{}, "", err
	}

	return user, token, nil
}

func GetUsers(userID int) ([]User, error) {
	var users []User
	query := `
		WITH last_messages AS (
			SELECT
				u.id AS user_id,
				u.first_name,
				u.last_name,
				u.nickname,
				u.email,
				u.age,
				u.gender,
				u.created_at as user_created_at,
				COALESCE(m.message, "") as last_message_content,
				COALESCE(m.sender, 0) as last_message_sender,
				COALESCE(strftime('%Y-%m-%dT%H:%M:%SZ', m.sent_at), "") AS sort_time
			FROM
				users u
			LEFT JOIN messages m
				ON m.id = (
					SELECT id
					FROM messages
					WHERE ((sender = u.id AND receiver = ? ) OR (sender = ? AND receiver = u.id))
					ORDER BY sent_at DESC
					LIMIT 1
				)
			WHERE
				u.id != ?
		)
		SELECT
			user_id AS id,
			first_name,
			last_name,
			nickname,
			email,
			age,
			gender,
			last_message_content,
			last_message_sender,
			sort_time
		FROM
			last_messages
		ORDER BY
			user_created_at, sort_time DESC;
`
	rows, err := DB.Query(query, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Nickname,
			&user.Email,
			&user.Age,
			&user.Gender,
			&user.LastMessage.Content,
			&user.LastMessage.SenderID,
			&user.LastMessage.SentAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetUserInfo(id int) (User, error) {
	var user User
	user.ID = id
	err := DB.QueryRow("SELECT first_name, last_name, nickname, email, age, gender FROM users WHERE id = ?", id).Scan(
		&user.FirstName,
		&user.LastName,
		&user.Nickname,
		&user.Email,
		&user.Age,
		&user.Gender)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUserPassword(user LoginRequest) (int, string, error) {
	var userID int
	var password string

	err := DB.QueryRow("SELECT id, password FROM users WHERE email = ? OR nickname = ?", user.Identifier, user.Identifier).Scan(&userID, &password)
	if err != nil {
		return 0, "", err
	}
	return userID, password, nil
}

func StoreNewUser(newUser RegistrationRequest, password string) (int64, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO users (first_name, last_name, nickname, email, age, gender, password) VALUES (?,?,?,?,?,?,?)`
	result, err := DB.Exec(query,
		newUser.FirstName,
		newUser.LastName,
		newUser.Nickname,
		newUser.Email,
		newUser.Age,
		newUser.Gender,
		hashedPassword)
	if err != nil {
		return 0, err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userId, nil
}
