package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"forum/server/models"
	"forum/server/utils"
	"forum/server/validators"
)

func IndexPosts(w http.ResponseWriter, r *http.Request) {
	statusCode, message, actionType, categoryID, page := validators.IndexPostsRequest(r)
	if statusCode != http.StatusOK {
		utils.JSONResponse(w, statusCode, message)
		return
	}
	// userID := 11

	var (
		posts []models.Post
		err   error
	)
	page *= 10

	switch actionType {
	case "index":
		posts, err = models.FetchPosts(page)
	case "category":
		posts, err = models.FetchPostsByCategory(categoryID, page)
	}

	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func ShowPost(w http.ResponseWriter, r *http.Request) {
	// Validate the request and extract the post ID
	statusCode, message, postID := validators.ShowPostRequest(r)
	if statusCode != http.StatusOK {
		utils.JSONResponse(w, statusCode, message)
		return
	}

	post, err := models.FetchPost(postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.JSONResponse(w, http.StatusNotFound, "Post not found")
		} else {
			log.Printf("Error fetching post with ID %d: %v", postID, err)
			utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	statusCode, message, title, content, categories := validators.CreatePostRequest(r)
	if statusCode != http.StatusOK {
		utils.JSONResponse(w, statusCode, message)
		return
	}

	user_id, _, valid := models.ValidSession(r)
	if !valid {
		utils.JSONResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := models.CheckCategories(categories)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid categories")
		return
	}

	post_id, err := models.StorePost(user_id, title, content)
	if err != nil {
		log.Println(err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for _, category := range categories {
		_, err = models.StorePostCategory(post_id, category)
		if err != nil {
			log.Println(err)
			utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	utils.JSONResponse(w, http.StatusOK, "success")
}
