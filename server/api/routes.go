package api

import (
	"net/http"

	"forum/server/controllers"
)

func Routes() http.Handler {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/", controllers.IndexAPIs)
	mux.HandleFunc("/posts", controllers.IndexPosts)
	mux.HandleFunc("/posts/{id}", controllers.ShowPost)
	mux.HandleFunc("/posts/create", controllers.CreatePost)
	mux.HandleFunc("/comment/create", controllers.CreateComment)
	mux.HandleFunc("/register", controllers.Register)
	mux.HandleFunc("/login", controllers.Login)
	mux.HandleFunc("/logout", controllers.Logout)

	return mux
}
