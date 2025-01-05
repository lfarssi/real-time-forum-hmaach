package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"

	"forum/server/utils"
)

func RecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// function that catch any panic
		defer func() {
			if err := recover(); err != nil {

				// log the panic and stack trace
				message := "Caught panic: %v, Stack trace: %s"
				log.Printf(message, err, string(debug.Stack()))
				utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		// continue processing the request and call the next handler
		next.ServeHTTP(w, r)
	}
}
