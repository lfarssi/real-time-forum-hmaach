package validators

import (
	"database/sql"
	"encoding/json"
	"html"
	"net/http"
	"strconv"
	"strings"

	"forum/server/models"
)

// validates a request for fetching comments by post ID.
// returns:
// - int: HTTP status code.
// - string: Error or success message.
// - int: postID
// - int: page of comments index
func GetCommentsRequest(r *http.Request) (int, string, int, int) {
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

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || postID < 1 {
		return http.StatusBadRequest, "Invalid post ID", 0, 0
	}

	return http.StatusOK, "success", postID, page
}

// Validates a request to create a comment.
// Returns:
// - models.CommentRequest: The parsed comment request.
// - int: HTTP status code.
// - string: Error or success message.
func CreateCommentRequest(r *http.Request) (models.CommentRequest, int, string) {
	if r.Method != http.MethodPost {
		return models.CommentRequest{}, http.StatusMethodNotAllowed, "Request method must be POST"
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return models.CommentRequest{}, http.StatusUnsupportedMediaType, "Content-Type must be 'application/json'"
	}

	var comment models.CommentRequest
	// Decode the JSON data
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		return models.CommentRequest{}, http.StatusBadRequest, "Invalid request body: unable to parse JSON data"
	}

	// Validate comment content
	comment.Content = html.EscapeString(strings.TrimSpace(comment.Content))
	if comment.Content == "" {
		return models.CommentRequest{}, http.StatusBadRequest, "Comment content cannot be empty"
	}
	if len(comment.Content) > 1800 {
		return models.CommentRequest{}, http.StatusBadRequest, "Comment content exceeds the maximum allowed length of 1800 characters"
	}

	// Validate Post ID
	if comment.PostID <= 0 {
		return models.CommentRequest{}, http.StatusBadRequest, "Post ID must be a positive integer"
	}

	// Check if the post exists
	if err := models.CheckPostExist(comment.PostID); err != nil {
		if err == sql.ErrNoRows {
			return models.CommentRequest{}, http.StatusBadRequest, "The specified post does not exist"
		}
		return models.CommentRequest{}, http.StatusInternalServerError, "An error occurred while verifying the post"
	}

	return comment, http.StatusOK, "Comment validated successfully"
}
