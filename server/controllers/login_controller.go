package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"forum/server/models"
	"forum/server/utils"
	"forum/server/validators"
)

func Login(w http.ResponseWriter, r *http.Request) {
	userRequest, statusCode, message := validators.LoginRequest(r)
	if statusCode != http.StatusOK {
		utils.JSONResponse(w, statusCode, message)
		return
	}

	// get user information from database
	user_id, hashedPassword, err := models.GetUserPassword(userRequest)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSONResponse(w, http.StatusNotFound, "User does not exist")
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Verify the password
	if match := utils.CheckPasswordHash(userRequest.Password, hashedPassword); !match {
		utils.JSONResponse(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	var userResponse models.User

	token, err := utils.GenerateToken()
	if err != nil {
		log.Println("Failed to create session: ", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}

	err = models.StoreSession(user_id, token, time.Now().Add(10*time.Hour))
	if err != nil {
		log.Println("Failed to store session into the database: ", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	userResponse, err = models.GetUserInfo(user_id)
	if err != nil {
		log.Println("Failed to fetch user's info: ", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	userResponse.Token = token

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResponse)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Context().Value("user_id").(int)
	err := models.DeleteUserSession(userID)
	if err != nil {
		log.Println("Error while logging out!")
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	utils.JSONResponse(w, http.StatusOK, "success")
}
