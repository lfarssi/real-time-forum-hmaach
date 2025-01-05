package controllers

import (
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

	_, err := models.StoreNewUser(user, password)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.nickname" {
			utils.JSONResponse(w, http.StatusNotAcceptable, "nickname already exists")
			return
		} else if err.Error() == "UNIQUE constraint failed: users.email" {
			utils.JSONResponse(w, http.StatusNotAcceptable, "email already exists")
			return
		}
		log.Println(err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	utils.JSONResponse(w, http.StatusOK, "success")
}
