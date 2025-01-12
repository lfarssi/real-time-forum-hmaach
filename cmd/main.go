package main

import (
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"forum/server/controllers"
	"forum/server/middlewares"
	"forum/server/models"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	// Parse the template during initialization
	_, err := template.ParseFiles("./web/index.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// Connect to the database
	err = models.Connect()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	err = models.CreateTables()
	if err != nil {
		log.Fatalf("error creating demo data: %v", err)
	}

	// err = models.CreateDemoData()
	// if err != nil {
	// 	log.Fatalf("error creating demo data: %v", err)
	// }
}

func routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", controllers.Index)
	mux.HandleFunc("/assets/", controllers.ServeStaticFiles)
	mux.HandleFunc("/api/register", controllers.Register)
	mux.HandleFunc("/api/login", controllers.Login)

	// Routes that require authentication
	mux.HandleFunc("/api/users", middlewares.IsAuth(controllers.IndexUsers))
	mux.HandleFunc("/api/posts", middlewares.IsAuth(controllers.IndexPosts))
	mux.HandleFunc("/api/posts/{id}/comments", middlewares.IsAuth(controllers.GetComments))
	mux.HandleFunc("/api/posts/create", middlewares.IsAuth(controllers.CreatePost))
	mux.HandleFunc("/api/posts/react", middlewares.IsAuth(controllers.ReactToPost))
	mux.HandleFunc("/api/comments/create", middlewares.IsAuth(controllers.CreateComment))
	mux.HandleFunc("/api/logout", middlewares.IsAuth(controllers.Logout))

	// WebSocket endpoint
	mux.HandleFunc("/ws", middlewares.IsAuthWebSocket(controllers.HandleWebSocket))

	// Create a rate limiter allowing 10 requests per minute
	rateLimiter := middlewares.NewRateLimiter(100, 1*time.Minute)

	// Apply RecoveryMiddleware, CORS, and RateLimiter globally
	return middlewares.CORS(middlewares.Recovery(rateLimiter.Middleware(mux)))
}

func main() {
	if len(os.Args) != 1 {
		log.Fatal("Too many arguments")
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: routes(),
	}

	log.Println("Server starting on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}
