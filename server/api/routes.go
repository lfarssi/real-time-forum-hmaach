package api

import (
	"net/http"
	"time"

	"forum/server/controllers"
	"forum/server/middlewares"
)

func Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", controllers.IndexAPIs)
	mux.HandleFunc("/posts", controllers.IndexPosts)
	mux.HandleFunc("/posts/{id}/comments", controllers.GetComments)
	mux.HandleFunc("/register", controllers.Register)
	mux.HandleFunc("/login", controllers.Login)

	// Routes that require authentication
	mux.HandleFunc("/posts/create", middlewares.IsAuth(controllers.CreatePost))
	mux.HandleFunc("/comments/create", middlewares.IsAuth(controllers.CreateComment))
	mux.HandleFunc("/logout", middlewares.IsAuth(controllers.Logout))

	// Create a rate limiter allowing 10 requests per minute
	rateLimiter := middlewares.NewRateLimiter(10, 1*time.Minute)

	// Apply RecoveryMiddleware, CORS, and RateLimiter globally
	return middlewares.CORS(middlewares.Recovery(rateLimiter.Middleware(mux)))
}
