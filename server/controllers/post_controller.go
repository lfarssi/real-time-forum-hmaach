package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"forum/server/models"
	"forum/server/utils"
	"forum/server/validators"
)

func IndexPosts(w http.ResponseWriter, r *http.Request) {
	statusCode, message, page := validators.IndexPostsRequest(r)
	if statusCode != http.StatusOK {
		utils.JSONResponse(w, statusCode, message)
		return
	}

	var (
		posts []models.Post
		err   error
	)
	limit := 10
	posts, err = models.FetchPosts(limit, page)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	statusCode, message, title, content, categories := validators.CreatePostRequest(r)
	if statusCode != http.StatusOK {
		utils.JSONResponse(w, statusCode, message)
		return
	}
	userID := r.Context().Value("user_id").(int)

	err := models.CheckCategories(categories)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid categories")
		return
	}

	post_id, err := models.StorePost(userID, title, content)
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
