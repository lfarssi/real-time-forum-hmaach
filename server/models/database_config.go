package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var (
	database     = "sqlite3"
	databasePath = "./server/database/database.db"
	schemaPath   = "./server/database/schema.sql"
	seedPath     = "./server/database/seed.sql"
	DB           *sql.DB
)

// connect to the database
func Connect() error {
	var err error
	DB, err = sql.Open(database, databasePath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}
	return nil
}

// CreateTables executes all queries from schema.sql
func CreateTables() error {
	content, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema.sql file: %v", err)
	}

	if _, err := DB.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to create tables %q: %v", string(content), err)
	}

	var catCount int
	if err = DB.QueryRow(`SELECT COUNT(*) FROM categories`).Scan(&catCount); err != nil {
		return fmt.Errorf("failed to get the count of categories: %v", err)
	}

	if catCount == 0 {
		query := `INSERT INTO categories (label) VALUES
            ('Technology'), ('Sport'),
            ('Business'),	('Health'),
            ('News'), ('devlopement');`

		if _, err = DB.Exec(query); err != nil {
			return fmt.Errorf("failed to insert categories into database: %v", err)
		}
	}
	return nil
}

// CreateDemoData generates and inserts fake data into the database
func CreateDemoData() error {
	content, err := os.ReadFile(seedPath)
	if err != nil {
		return fmt.Errorf("failed to read seed.sql file: %v", err)
	}

	if _, err = DB.Exec(string(content)); err != nil {
		log.Printf("failed to insert demo data %q: %v\n", string(content), err)
		return err
	}

	return nil
}
