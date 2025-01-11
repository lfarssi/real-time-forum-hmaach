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
	post, statusCode, message := validators.CreatePostRequest(r)
	if statusCode != http.StatusOK {
		utils.JSONResponse(w, statusCode, message)
		return
	}

	post.UserID = r.Context().Value("user_id").(int)


	postID, err := models.StorePost(post)
	if err != nil {
		log.Println(err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	for _, category := range post.Categories {
		_, err = models.StorePostCategory(postID, category)
		if err != nil {
			log.Println(err)
			utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	utils.JSONResponse(w, http.StatusOK, "success")
}

func ReactToPost(w http.ResponseWriter, r *http.Request) {
	reaction, status, message := validators.ReactToPostRequest(r)
	if status != http.StatusOK {
		utils.JSONResponse(w, status, message)
		return
	}
	reaction.UserID = r.Context().Value("user_id").(int)
	err := models.ReactToPost(reaction)
	if err != nil {
		log.Println(err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	utils.JSONResponse(w, http.StatusOK, "success")
}
