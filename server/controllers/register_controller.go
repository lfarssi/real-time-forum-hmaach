package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"forum/server/models"
	"forum/server/utils"
	"forum/server/validators"
)

func Register(w http.ResponseWriter, r *http.Request) {
	user, password, statusCode, message := validators.RegisterRequest(r)
	if statusCode != http.StatusOK {
		utils.JSONResponse(w, statusCode, message)
		return
	}
	
	userID, err := models.StoreNewUser(user, password)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.nickname" {
			utils.JSONResponse(w, http.StatusNotAcceptable, "nickname already exists")
			return
		} else if err.Error() == "UNIQUE constraint failed: users.email" {
			utils.JSONResponse(w, http.StatusNotAcceptable, "email already exists")
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	userResponse, token, err := models.GenerateSession(int(userID))
	if err != nil {
		log.Println(err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	broadcastOnlineUserList()
	broadcastMessage("refresh-users")
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "success", "user": userResponse, "token": token, "status": 200})
}
