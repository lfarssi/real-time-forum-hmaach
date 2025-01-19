package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

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
	userID, hashedPassword, err := models.GetUserPassword(userRequest)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSONResponse(w, http.StatusUnauthorized, "Invalid nickname or email")
			return
		}
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Verify the password
	if !utils.ComparePassword(userRequest.Password, hashedPassword) {
		utils.JSONResponse(w, http.StatusUnauthorized, "incorrect password")
		return
	}
	userResponse, token, err := models.GenerateSession(userID)
	if err != nil {
		log.Println(err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "success", "user": userResponse, "token": token, "status": 200})
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
