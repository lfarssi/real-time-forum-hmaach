package validators

import (
	"encoding/json"
	"html"
	"net/http"
	"strings"

	"forum/server/models"
	"forum/server/utils"
)

// Validates a registration request.
// Returns:
// - user: user information for registration
// - string: password for registration
// - int: HTTP status code.
// - string: Error or success message.
func RegisterRequest(r *http.Request) (models.RegistrationRequest, string, int, string) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		return models.RegistrationRequest{}, "", http.StatusMethodNotAllowed, "Invalid HTTP method"
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return models.RegistrationRequest{}, "", http.StatusUnsupportedMediaType, "Content-Type must be application/json"
	}

	var user models.RegistrationRequest

	// Decode the JSON data
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Invalid JSON data"
	}

	// Sanitize inputs
	user.FirstName = html.EscapeString(strings.TrimSpace(user.FirstName))
	user.LastName = html.EscapeString(strings.TrimSpace(user.LastName))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.Nickname = html.EscapeString(strings.TrimSpace(user.Nickname))
	user.Gender = html.EscapeString(strings.TrimSpace(user.Gender))
	password := html.EscapeString(strings.TrimSpace(user.Password))
	passwordConfirmation := html.EscapeString(strings.TrimSpace(user.PasswordConfirmation))

	// Validate First Name
	if len(user.FirstName) < 2 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "First name must be at least 2 characters long"
	}
	if len(user.FirstName) > 50 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "First name cannot exceed 50 characters"
	}
	if strings.ContainsAny(user.FirstName, "1234567890") {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "First name cannot contain numbers"
	}

	// Validate Last Name
	if len(user.LastName) < 2 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Last name must be at least 2 characters long"
	}
	if len(user.LastName) > 50 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Last name cannot exceed 50 characters"
	}
	if strings.ContainsAny(user.LastName, "1234567890") {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Last name cannot contain numbers"
	}

	// Validate email
	if len(user.Email) > 100 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Email cannot exceed 100 characters"
	}
	if user.Email == "" {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Email is required"
	}
	if !utils.IsValidEmail(user.Email) {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Invalid email format"
	}

	// Validate nickname
	if len(user.Nickname) < 4 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Nickname must be at least 4 characters long"
	}
	if len(user.Nickname) > 20 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Nickname cannot exceed 20 characters"
	}
	if strings.Contains(user.Nickname, " ") {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Nickname cannot contain spaces"
	}

	// Validate gender
	if len(user.Gender) > 10 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Gender cannot exceed 10 characters"
	}
	if user.Gender != "male" && user.Gender != "female" && user.Gender != "other" {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Invalid gender. Must be 'male', 'female', or 'other'"
	}

	// Validate age
	if user.Age <= 0 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Invalid Age value"
	}
	if user.Age > 120 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Age cannot exceed 120"
	}

	// Validate password
	if password != passwordConfirmation {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Passwords do not match"
	}
	if len(password) < 6 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Password must be at least 6 characters long"
	}
	if len(password) > 128 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Password cannot exceed 128 characters"
	}
	if !utils.ContainsSpecialChar(password) {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Password must contain at least one special character"
	}
	if !utils.ContainsUppercase(password) {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Password must contain at least one uppercase letter"
	}
	if !utils.ContainsLowercase(password) {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Password must contain at least one lowercase letter"
	}
	if !utils.ContainsNumber(password) {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Password must contain at least one number"
	}

	return user, password, http.StatusOK, "success"
}
