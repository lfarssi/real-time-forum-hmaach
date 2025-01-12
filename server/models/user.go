package models

import (
	"forum/server/utils"
)

// represents the data for user registration.
type RegistrationRequest struct {
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Email                string `json:"email"`
	Nickname             string `json:"nickname"`
	Gender               string `json:"gender"`
	Age                  int    `json:"age"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

// represents the data for user login.
type LoginRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
}

func GetUsers(userID int) ([]User, error) {
	var users []User
	query := `
		SELECT 	
			id, first_name, last_name, nickname, email, age, gender 
		FROM 
			users 
		WHERE NOT id = ?
		ORDER BY 
			created_at DESC`
	rows, err := DB.Query(query, userID)
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
			&user.Gender)
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
	var (
		userID   int
		password string
		err      error
	)
	if user.Email != "" {
		err = DB.QueryRow("SELECT id, password FROM users WHERE email = ?", user.Email).Scan(&userID, &password)
	} else {
		err = DB.QueryRow("SELECT id, password FROM users WHERE nickname =?", user.Nickname).Scan(&userID, &password)
	}
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

	query := `INSERT INTO users (first_name, last_name, nickname, email,age, gender, password) VALUES (?,?,?,?,?,?,?)`
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

	userID, _ := result.LastInsertId()

	return userID, nil
}
