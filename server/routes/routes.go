package routes

import (
	"net/http"
	"time"

	"forum/server/controllers"
	"forum/server/middlewares"
)

func Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", controllers.Index)
	mux.HandleFunc("/assets/", controllers.ServeStaticFiles)
	mux.HandleFunc("/api/register", controllers.Register)
	mux.HandleFunc("/api/login", controllers.Login)

	// Routes that require authentication
	mux.HandleFunc("/api/posts", middlewares.IsAuth(controllers.IndexPosts))
	mux.HandleFunc("/api/posts/{id}/comments", middlewares.IsAuth(controllers.GetComments))
	mux.HandleFunc("/api/posts/create", middlewares.IsAuth(controllers.CreatePost))
	mux.HandleFunc("/api/posts/react", middlewares.IsAuth(controllers.ReactToPost))
	mux.HandleFunc("/api/comments/create", middlewares.IsAuth(controllers.CreateComment))
	mux.HandleFunc("/api/logout", middlewares.IsAuth(controllers.Logout))

	// Create a rate limiter allowing 10 requests per minute
	rateLimiter := middlewares.NewRateLimiter(10, 1*time.Minute)

	// Apply RecoveryMiddleware, CORS, and RateLimiter globally
	return middlewares.CORS(middlewares.Recovery(rateLimiter.Middleware(mux)))
}
