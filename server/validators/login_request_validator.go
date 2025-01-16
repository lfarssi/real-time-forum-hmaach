package validators

import (
	"encoding/json"
	"html"
	"net/http"
	"strings"

	"forum/server/models"
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
		return models.LoginRequest{}, http.StatusUnsupportedMediaType, "Content-Type must be application/json"
	}

	var user models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return models.LoginRequest{}, http.StatusBadRequest, "Invalid JSON data"
	}

	user.Identifier = html.EscapeString(strings.TrimSpace(user.Identifier))
	user.Password = html.EscapeString(user.Password)

	if user.Identifier == "" {
		return models.LoginRequest{}, http.StatusBadRequest, "identifier is required"
	}
	if user.Password == "" {
		return models.LoginRequest{}, http.StatusBadRequest, "Password is required"
	}
	if len(user.Password) > 50 {
		return models.LoginRequest{}, http.StatusBadRequest, "Password cannot exceed 50 characters"
	}

	return user, http.StatusOK, "success"
}
