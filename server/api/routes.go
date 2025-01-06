package api

import (
	"net/http"
	"time"

	"forum/server/controllers"
	"forum/server/middlewares"
)

func Routes() http.Handler {
	mux := http.NewServeMux()

	// Create a rate limiter allowing 10 requests per minute
	rateLimiter := middlewares.NewRateLimiter(10, 1*time.Minute)

	mux.HandleFunc("/", middlewares.RecoveryMiddleware(rateLimiter.Middleware(controllers.IndexAPIs)))
	mux.HandleFunc("/posts", middlewares.RecoveryMiddleware(rateLimiter.Middleware(controllers.IndexPosts)))
	mux.HandleFunc("/posts/{id}/comments", middlewares.RecoveryMiddleware(rateLimiter.Middleware(controllers.GetComments)))
	mux.HandleFunc("/register", middlewares.RecoveryMiddleware(rateLimiter.Middleware(controllers.Register)))
	mux.HandleFunc("/login", middlewares.RecoveryMiddleware(rateLimiter.Middleware(controllers.Login)))

	// routes that require authentication
	mux.HandleFunc("/posts/create", middlewares.RecoveryMiddleware(middlewares.AuthMiddleware(rateLimiter.Middleware(controllers.CreatePost))))
	mux.HandleFunc("/comment/create", middlewares.RecoveryMiddleware(middlewares.AuthMiddleware(rateLimiter.Middleware(controllers.CreateComment))))
	mux.HandleFunc("/logout", middlewares.RecoveryMiddleware(middlewares.AuthMiddleware(rateLimiter.Middleware(controllers.Logout))))

	return mux
}
