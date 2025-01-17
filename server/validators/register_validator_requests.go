package validators

import (
	"encoding/json"
	"html"
	"net/http"
	"regexp"
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
	if r.Method != http.MethodPost {
		return models.RegistrationRequest{}, "", http.StatusMethodNotAllowed, "Invalid HTTP method"
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return models.RegistrationRequest{}, "", http.StatusUnsupportedMediaType, "Content-Type must be application/json"
	}

	var user models.RegistrationRequest
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

	// Validate First Name
	if len(user.FirstName) < 3 && len(user.FirstName) > 20 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "First name must be between 3 and 20 characters long"
	}
	if valid, err := regexp.MatchString(`^[a-zA-Z]+$`, user.FirstName); err != nil || !valid {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "First name can only contain characters"
	}

	// Validate Last Name
	if len(user.LastName) < 3 && len(user.LastName) > 20 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Last name must be between 3 and 20 characters long"
	}
	if valid, err := regexp.MatchString(`(?i)^[a-z]+$`, user.LastName); err != nil || !valid {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Last name can only contain characters"
	}

	// Validate email
	if len(user.Email) < 5 && len(user.Email) > 50 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Email must be between 5 and 50 characters long"
	}
	if isValid, err := regexp.MatchString(`(?i)^[a-z0-9]+\.?[a-z0-9]+@[a-z0-9]+\.[a-z]+$`, user.Email); err != nil || !isValid {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Invalid email format"
	}

	// Validate nickname
	if len(user.Nickname) < 3 && len(user.Nickname) > 20 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Nickname must be between 3 and 20 characters long"
	}
	if isValid, err := regexp.MatchString(`^[a-z0-9]+$`, user.Nickname); err != nil || !isValid {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Nickname must be in lowercase letter"
	}

	// Validate gender
	if len(user.Gender) == 0 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Gender is required"
	}
	if user.Gender != "male" && user.Gender != "female" {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Invalid gender. Must be 'male' or 'female'"
	}

	// Validate age
	if user.Age < 18 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Our policy requires age to be bigger than 18"
	}
	if user.Age > 120 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Age cannot exceed 120"
	}

	// Validate password
	if len(password) < 6 && len(password) > 50 {
		return models.RegistrationRequest{}, "", http.StatusBadRequest, "Password must be between 6 and 50 characters long"
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
