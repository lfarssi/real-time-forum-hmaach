package validators

import (
	"net/http"
	"strconv"
	"strings"
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
// - int: HTTP status code.
// - string: Error or success message.
// - string: Comment content.
// - int: Post ID.
func CreateCommentRequest(r *http.Request) (int, string, string, int) {
	if r.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, "Invalid HTTP method", "", 0
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return http.StatusUnsupportedMediaType, "Content-Type must be application/json", "", 0
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		return http.StatusBadRequest, "Failed to parse form data", "", 0
	}

	// Validate comment content
	content := strings.TrimSpace(r.FormValue("comment"))
	if content == "" {
		return http.StatusBadRequest, "Comment content cannot be empty", "", 0
	}
	if len(content) > 1800 {
		return http.StatusBadRequest, "Comment content exceeds the maximum length of 1800 characters", "", 0
	}

	// Validate Post ID
	postIDStr := r.FormValue("postid")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		return http.StatusBadRequest, "Post ID must be a valid positive integer", "", 0
	}

	return http.StatusOK, "success", content, postID
}
