package middlewares

import (
	"context"
	"net/http"

	"forum/server/models"
	"forum/server/utils"
)

// IsAuth combined function for both HTTP and WebSocket requests
func IsAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string

		// Check if it's a WebSocket connection (Upgrade header check)
		if r.Header.Get("Upgrade") == "websocket" {
			token = r.URL.Query().Get("token")
		} else {
			token = r.Header.Get("Authorization")
		}

		if token == "" {
			utils.JSONResponse(w, http.StatusUnauthorized, "Token is required")
			return
		}

		userID, isValid, message := models.ValidSession(token)
		if !isValid {
			if message == "Internal Server Error" {
				utils.JSONResponse(w, http.StatusInternalServerError, message)
			} else {
				utils.JSONResponse(w, http.StatusUnauthorized, message)
			}
			return
		}

		// Attach userID to request context for further use
		r = r.WithContext(context.WithValue(r.Context(), "user_id", userID))

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	}
}
