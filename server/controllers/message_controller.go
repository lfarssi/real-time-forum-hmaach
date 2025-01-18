package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"forum/server/models"
	"forum/server/utils"
	"forum/server/validators"
)

func GetConvertation(w http.ResponseWriter, r *http.Request) {
	statusCode, message, senderID, page := validators.GetConvertationRequest(r)
	if statusCode != http.StatusOK {
		utils.JSONResponse(w, statusCode, message)
		return
	}

	sender, err := models.GetUserInfo(senderID)
	if err != nil {
		log.Println("Failed to fetch sender's info: ", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	userID := r.Context().Value("user_id").(int)
	limit := 10

	if senderID == userID {
		utils.JSONResponse(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Fetch messages from the database and return them as JSON
	messages, err := models.GetMessages(userID, senderID, limit, page)
	if err != nil {
		log.Println("Failed to fetch messages: ", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"sender": sender, "messages": messages})
}
