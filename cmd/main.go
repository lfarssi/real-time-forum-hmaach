package main

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"forum/server/models"
	"forum/server/routes"

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
}

func main() {
	if len(os.Args) != 1 {
		log.Fatal("Too many arguments")
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: routes.Routes(),
	}

	log.Println("Server starting on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}
