package controllers

import (
	"net/http"

	"forum/server/models"
	"forum/server/utils"
	"forum/server/validators"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	statusCode, message, content, postID := validators.CreateCommentRequest(r)
	if statusCode != http.StatusOK {
		utils.JSONResponse(w, statusCode, message)
		return
	}
	userID := r.Context().Value("user_id").(int)

	// Store the comment using the models package
	_, err := models.StoreComment(userID, postID, content)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	utils.JSONResponse(w, http.StatusOK, "success")
}
