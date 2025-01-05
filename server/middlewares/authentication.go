package middlewares

import (
	"context"
	"net/http"

	"forum/server/models"
	"forum/server/utils"
)

// AuthMiddleware checks if the user is authenticated by validating the session token.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		token := r.Header.Get("Authorization")
		if token == "" {
			utils.JSONResponse(w, http.StatusUnauthorized, "Token is required")
			return
		}

		// Validate the session token
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
