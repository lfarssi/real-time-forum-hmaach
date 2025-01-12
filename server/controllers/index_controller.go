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
	// Check if the requested path is not the root path
	if r.URL.Path != "/" {
		utils.JSONResponse(w, http.StatusNotFound, "Page Not Found")
		return
	}

	// Check if the request method is GET
	if r.Method != http.MethodGet {
		utils.JSONResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// Execute the pre-parsed template
	t, _ := template.ParseFiles("./web/index.html")
	err := t.Execute(w, nil)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

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

	var connectedIDs []int
	for id := range ConnectedUsers {
		connectedIDs = append(connectedIDs, id)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"users": users, "connected": connectedIDs})
}

// ServeStaticFiles returns a handler function for serving static files
func ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	// Get clean file path and prevent directory traversal
	filePath := filepath.Clean("./web" + r.URL.Path)

	// block access to dirictories
	if info, err := os.Stat(filePath); err != nil || info.IsDir() {
		utils.JSONResponse(w, http.StatusNotFound, "Page not found")
		return
	}

	// Serve the file
	http.ServeFile(w, r, filePath)
}
