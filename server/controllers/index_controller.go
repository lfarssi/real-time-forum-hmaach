package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"forum/server/models"
	"forum/server/utils"
)

// Index handles the root route and serves the index.html template
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// Execute the pre-parsed template
	t, _ := template.ParseFiles("./web/index.html")
	if err := t.Execute(w, nil); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

// IndexUsers handles the root route and serves the index.html template
func IndexUsers(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// Fetch users from the database and return them as JSON
	userID := r.Context().Value("user_id").(int)
	users, err := models.GetUsers(userID)
	if err != nil {
		log.Println("Failed to fetch users: ", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"users": users, "status": http.StatusOK})
}

// ServeStaticFiles returns a handler function for serving static files
func ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	filePath := filepath.Clean("./web" + r.URL.Path)

	if info, err := os.Stat(filePath); err != nil || info.IsDir() {
		utils.JSONResponse(w, http.StatusNotFound, "Page not found")
		return
	}

	http.ServeFile(w, r, filePath)
}
