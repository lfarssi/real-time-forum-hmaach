package validators

import (
	"encoding/json"
	"html"
	"net/http"
	"strconv"
	"strings"

	"forum/server/models"
)

// validates a request for fetching a conversation.
// returns:
// - int: HTTP status code.
// - string: Error or success message.
// - int: page of messages to fetch
// - int: sender ID
func GetConvertationRequest(r *http.Request) (int, string, int, int) {
	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, "Invalid HTTP method", 0, 0
	}

	err := r.ParseForm()
	if err != nil {
		return http.StatusBadRequest, "Failed to parse form data", 0, 0
	}

	page := 0
	pageStr := r.FormValue("page")
	if pageStr != "" {
		page, err = strconv.Atoi(r.FormValue("page"))
		if err != nil || page < 1 {
			return http.StatusBadRequest, "Invalid page number", 0, 0
		}
		page-- // in the databse the page number should start from 0
	}

	// validate message sender id
	senderID, err := strconv.Atoi(r.PathValue("sender"))
	if err != nil || senderID < 0 {
		return http.StatusBadRequest, "Invalid sender ID", 0, 0
	}
	_, err = models.GetUserInfo(senderID)
	if err != nil {
		return http.StatusBadRequest, "Invalid sender ID", 0, 0
	}

	return http.StatusOK, "success", senderID, page
}

// validates a request to create a new post.
// Returns:
// - models.MessageRequest: The validated post request structure.
// - int: HTTP status code.
// - string: Error or success message.
func SendMessageRequest(r *http.Request) (models.MessageRequest, int, string) {
	// Check HTTP method
	if r.Method != http.MethodPost {
		return models.MessageRequest{}, http.StatusMethodNotAllowed, "Only POST method is allowed"
	}

	// Check Content-Type header
	if r.Header.Get("Content-Type") != "application/json" {
		return models.MessageRequest{}, http.StatusUnsupportedMediaType, "Content-Type must be 'application/json'"
	}

	var message models.MessageRequest

	// Decode the JSON data
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		return models.MessageRequest{}, http.StatusBadRequest, "Invalid JSON data: unable to parse request body"
	}

	// Sanitize and validate title
	message.Text = html.EscapeString(strings.TrimSpace(message.Text))
	if message.Text == "" {
		return models.MessageRequest{}, http.StatusBadRequest, "The title field is required and cannot be empty"
	}
	if len(message.Text) > 1000 {
		return models.MessageRequest{}, http.StatusBadRequest, "The message must not exceed 1000 characters"
	}

	// validate message sender id
	_, err := models.GetUserInfo(message.Sender)
	if err != nil {
		return models.MessageRequest{}, http.StatusBadRequest, "Invalid sender ID"
	}

	return message, http.StatusOK, "success"
}
