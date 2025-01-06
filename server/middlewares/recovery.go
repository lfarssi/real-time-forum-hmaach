package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"

	"forum/server/utils"
)

// RecoveryMiddleware wraps an http.Handler to recover from panics
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// log the panic and stack trace
				message := "Caught panic: %v, Stack trace: %s"
				log.Printf(message, err, string(debug.Stack()))
				utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
