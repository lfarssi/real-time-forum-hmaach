package main

import (
	"log"
	"net/http"
	"os"

	"forum/server/api"
	"forum/server/models"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	// Connect to the database
	err := models.Connect()
	if err != nil {
		log.Panic("Database connection error:", err)
	}
	err = models.CreateTables()
	if err != nil {
		log.Fatal("error creating demo data:", err)
	}
}

func main() {
	if len(os.Args) != 1 {
		log.Fatal("Too many arguments")
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: api.Routes(),
	}

	log.Println("Server starting on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}
