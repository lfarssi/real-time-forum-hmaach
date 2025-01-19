package validators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/server/models"
)

func ChatMessageRequest(userID int, data []byte) (models.Message, error) {
	var message models.Message

	err := json.NewDecoder(bytes.NewReader(data)).Decode(&message)
	if err != nil {
		return models.Message{}, fmt.Errorf("invalid JSON data: unable to parse request body")
	}

	if userID == message.ReceiverID {
		return models.Message{}, fmt.Errorf("invalid receiver ID")
	}

	// validate message content
	if len(message.Content) == 0 {
		return models.Message{}, fmt.Errorf("message content cannot be empty")
	}
	if len(message.Content) > 200 {
		return models.Message{}, fmt.Errorf("message must not exceed 200 characters")
	}
	messageLines := strings.Count(message.Content, "\n") + 1
	if messageLines > 3 {
		return models.Message{}, fmt.Errorf("message cannot exceed 3 lines")
	}

	// validate message receiver id
	if message.ReceiverID < 0 {
		return models.Message{}, fmt.Errorf("invalid receiver ID")
	}
	_, err = models.GetUserInfo(message.ReceiverID)
	if err != nil {
		return models.Message{}, fmt.Errorf("invalid receiver ID")
	}

	message.Type = "message"

	return message, nil
}

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
