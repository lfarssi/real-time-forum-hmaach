package api

import (
	"net/http"

	"forum/server/controllers"
	"forum/server/middlewares"
)

func Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", middlewares.RecoveryMiddleware(controllers.IndexAPIs))
	mux.HandleFunc("/posts", middlewares.RecoveryMiddleware(controllers.IndexPosts))
	mux.HandleFunc("/posts/{id}", middlewares.RecoveryMiddleware(controllers.ShowPost))
	mux.HandleFunc("/register", middlewares.RecoveryMiddleware(controllers.Register))
	mux.HandleFunc("/login", middlewares.RecoveryMiddleware(controllers.Login))

	// routes that require authentication
	mux.HandleFunc("/posts/create", middlewares.RecoveryMiddleware(middlewares.AuthMiddleware(controllers.CreatePost)))
	mux.HandleFunc("/comment/create", middlewares.RecoveryMiddleware(middlewares.AuthMiddleware(controllers.CreateComment)))
	mux.HandleFunc("/logout", middlewares.RecoveryMiddleware(middlewares.AuthMiddleware(controllers.Logout)))

	return mux
}
