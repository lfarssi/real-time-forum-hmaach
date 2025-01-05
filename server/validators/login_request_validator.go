package validators

import (
	"encoding/json"
	"html"
	"net/http"
	"strings"

	"forum/server/models"
	"forum/server/utils"
)

// Validates a login request.
// Returns:
// - user: user information for login
// - int: HTTP status code.
// - string: Error or success message.
func LoginRequest(r *http.Request) (models.LoginRequest, int, string) {
	if r.Method != http.MethodPost {
		return models.LoginRequest{}, http.StatusMethodNotAllowed, "Invalid HTTP method"
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return models.LoginRequest{}, http.StatusBadRequest, "Content-Type must be application/json"
	}

	var user models.LoginRequest

	// Decode the JSON data
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return models.LoginRequest{}, http.StatusBadRequest, "Invalid JSON data"
	}

	user.Nickname = html.EscapeString(strings.TrimSpace(user.Nickname))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	password := html.EscapeString(strings.TrimSpace(user.Password))

	if (user.Nickname == "" && user.Email == "") || (user.Nickname != "" && user.Email != "") {
		return models.LoginRequest{}, http.StatusBadRequest, "Either Nickname or Email is required"
	}

	// Validate Nickname
	if user.Nickname != "" && len(user.Nickname) > 20 {
		return models.LoginRequest{}, http.StatusBadRequest, "Nickname cannot exceed 20 characters"
	}

	// Validate Email (if provided)
	if user.Email != "" {
		if !utils.IsValidEmail(user.Email) || len(user.Nickname) > 100 {
			return models.LoginRequest{}, http.StatusBadRequest, "Invalid email format"
		}
	}

	// Validate Password
	if user.Password == "" {
		return models.LoginRequest{}, http.StatusBadRequest, "Password is required"
	}
	if len(password) > 128 {
		return models.LoginRequest{}, http.StatusBadRequest, "Password cannot exceed 128 characters"
	}

	return user, http.StatusOK, "success"
}
