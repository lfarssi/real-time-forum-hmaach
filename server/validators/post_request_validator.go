package validators

import (
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
	if len(post.Title) > 100 {
		return models.PostRequest{}, http.StatusBadRequest, "The title must not exceed 100 characters"
	}

	// Sanitize and validate content
	post.Content = html.EscapeString(strings.TrimSpace(post.Content))
	if post.Content == "" {
		return models.PostRequest{}, http.StatusBadRequest, "The content field is required and cannot be empty"
	}
	if len(post.Content) > 3000 {
		return models.PostRequest{}, http.StatusBadRequest, "The content must not exceed 3000 characters"
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

	return post, http.StatusOK, "Post request validated successfully"
}
