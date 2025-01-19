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

// validates a request for posts index.
// returns:
// - int: HTTP status code.
// - string: Error or success message.
// - int: page of posts index
func IndexPostsRequest(r *http.Request) (int, string, int) {
	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, "Invalid HTTP method", 0
	}

	err := r.ParseForm()
	if err != nil {
		return http.StatusBadRequest, "Failed to parse form data", 0
	}

	page := 0
	pageStr := r.FormValue("page")
	if pageStr != "" {
		page, err = strconv.Atoi(r.FormValue("page"))
		if err != nil || page < 1 {
			return http.StatusBadRequest, "Invalid page number", 0
		}
		page-- // in the databse the page number should start from 0
	}

	return http.StatusOK, "success", page
}

// validates a request to create a new post.
// Returns:
// - models.PostRequest: The validated post request structure.
// - int: HTTP status code.
// - string: Error or success message.
func CreatePostRequest(r *http.Request) (models.PostRequest, int, string) {
	// Check HTTP method
	if r.Method != http.MethodPost {
		return models.PostRequest{}, http.StatusMethodNotAllowed, "Only POST method is allowed"
	}

	// Check Content-Type header
	if r.Header.Get("Content-Type") != "application/json" {
		return models.PostRequest{}, http.StatusUnsupportedMediaType, "Content-Type must be 'application/json'"
	}

	var post models.PostRequest

	// Decode the JSON data
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		return models.PostRequest{}, http.StatusBadRequest, "Invalid JSON data: unable to parse request body"
	}

	// Sanitize and validate title
	post.Title = html.EscapeString(strings.TrimSpace(post.Title))
	if post.Title == "" {
		return models.PostRequest{}, http.StatusBadRequest, "The title field is required and cannot be empty"
	}
	if len(post.Title) > 70 {
		return models.PostRequest{}, http.StatusBadRequest, "The title must not exceed 70 characters"
	}

	// Sanitize and validate content
	post.Content = html.EscapeString(strings.TrimSpace(post.Content))
	if post.Content == "" {
		return models.PostRequest{}, http.StatusBadRequest, "The content field is required and cannot be empty"
	}
	if len(post.Content) > 1000 {
		return models.PostRequest{}, http.StatusBadRequest, "The content must not exceed 1000 characters"
	}
	contentLines := strings.Count(post.Content, "\n") + 1
    if contentLines > 5 {
        return models.PostRequest{}, http.StatusBadRequest, "Content cannot exceed 5 lines"
    }

	// Validate categories
	if len(post.Categories) == 0 {
		return models.PostRequest{}, http.StatusBadRequest, "At least one category must be selected"
	}
	for _, cat := range post.Categories {
		if cat < 1 {
			return models.PostRequest{}, http.StatusBadRequest, "Category IDs must be positive integers"
		}
	}

	// Check if categories exist in the database
	if err := models.CheckCategoriesExist(post.Categories); err != nil {
		return models.PostRequest{}, http.StatusBadRequest, "One or more category IDs are invalid"
	}

	return post, http.StatusOK, "success"
}

// validates a request to react to a post.
// Returns:
// - models.Reaction: The validated reaction request structure.
// - int: HTTP status code.
// - string: Error or success message.
func ReactToPostRequest(r *http.Request) (models.Reaction, int, string) {
	// Check HTTP method
	if r.Method != http.MethodPost {
		return models.Reaction{}, http.StatusMethodNotAllowed, "Only POST method is allowed"
	}

	// Check Content-Type header
	if r.Header.Get("Content-Type") != "application/json" {
		return models.Reaction{}, http.StatusUnsupportedMediaType, "Content-Type must be 'application/json'"
	}

	var reaction models.Reaction

	// Decode the JSON data
	if err := json.NewDecoder(r.Body).Decode(&reaction); err != nil {
		return models.Reaction{}, http.StatusBadRequest, "Invalid JSON data: unable to parse request body"
	}

	// Sanitize and validate title
	reaction.Type = html.EscapeString(strings.TrimSpace(reaction.Type))
	if reaction.Type == "" || (reaction.Type != "like" && reaction.Type != "dislike") {
		return models.Reaction{}, http.StatusBadRequest, "Invalid reaction type: must be 'like' or 'dislike'"
	}

	// Validate Post ID
	if reaction.PostID <= 0 {
		return models.Reaction{}, http.StatusBadRequest, "Post ID must be a positive integer"
	}

	// Check if the post exists
	if err := models.CheckPostExist(reaction.PostID); err != nil {
		if err == sql.ErrNoRows {
			return models.Reaction{}, http.StatusBadRequest, "The specified post does not exist"
		}
		return models.Reaction{}, http.StatusInternalServerError, "An error occurred while verifying the post"
	}

	return reaction, http.StatusOK, "success"
}
