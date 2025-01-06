package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"forum/server/models"
	"forum/server/utils"
	"forum/server/validators"
)

func GetComments(w http.ResponseWriter, r *http.Request) {
	statusCode, message, postID, page := validators.GetCommentsRequest(r)
	if statusCode != http.StatusOK {
		utils.JSONResponse(w, statusCode, message)
		return
	}

	var (
		comments []models.Comment
		err      error
	)
	limit := 10

	comments, err = models.FetchCommentsByPostID(postID, limit, page)
	if err != nil {
		log.Println(err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	comment, statusCode, message := validators.CreateCommentRequest(r)
	if statusCode != http.StatusOK {
		utils.JSONResponse(w, statusCode, message)
		return
	}
	comment.UserID = r.Context().Value("user_id").(int)

	// Store the comment using the models package
	_, err := models.StoreComment(comment)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	utils.JSONResponse(w, http.StatusOK, "success")
}
