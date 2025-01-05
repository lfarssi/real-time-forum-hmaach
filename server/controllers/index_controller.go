package controllers

import (
	"encoding/json"
	"net/http"

	"forum/server/config"
)

func IndexAPIs(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config.APIs)
}
