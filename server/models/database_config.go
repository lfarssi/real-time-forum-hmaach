package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"forum/server/utils"
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
	if err = DB.QueryRow(`SELECT COUNT(id) FROM categories`).Scan(&catCount); err != nil {
		return fmt.Errorf("failed to get the count of categories: %v", err)
	}

	if catCount == 0 {
		query := `INSERT INTO categories (label) VALUES
            ('Technology'), ('Sport'),
            ('Business'),	('Health'),
            ('News');`

		if _, err = DB.Exec(query); err != nil {
			return fmt.Errorf("failed to insert categories into database: %v", err)
		}
	}
	return nil
}

// CreateDemoData generates and inserts fake data into the database
func CreateDemoData() error {
	users := []RegistrationRequest{
		{"Hamza", "Maach", "hmaach", "hmaach@email.com", 25, "male", "123456789"},
		{"Mohammed", "Elfarssi", "melfarss", "melfarss@email.com", 28, "male", "123456789"},
		{"Fahd", "Agnouz", "fagnou", "fagnou@email.com", 30, "male", "123456789"},
		{"Yassine", "Elmach", "yelmach", "yelmach@email.com", 22, "male", "123456789"},
		{"Karim", "El Fassi", "karim_fassi", "karim.elfassi@email.com", 35, "male", "123456789"},
		{"Nadia", "Ghazali", "nadia_gh", "nadia.ghazali@email.com", 27, "female", "123456789"},
		{"Mehdi", "Houssaini", "mehdi_h", "mehdi.houssaini@email.com", 31, "male", "123456789"},
		{"Leila", "Idrissi", "leila_id", "leila.idrissi@email.com", 29, "female", "123456789"},
		{"Omar", "Jalal", "omar_j", "omar.jalal@email.com", 33, "male", "123456789"},
		{"Sanaa", "Kadiri", "sanaa_k", "sanaa.kadiri@email.com", 26, "female", "123456789"},
		{"Rachid", "Lahlou", "rachid_l", "rachid.lahlou@email.com", 34, "male", "123456789"},
		{"Salma", "Mansouri", "salma_m", "salma.mansouri@email.com", 24, "female", "123456789"},
		{"Hamza", "Najjar", "hamza_n", "hamza.najjar@email.com", 28, "male", "123456789"},
		{"Zineb", "Ouazzani", "zineb_o", "zineb.ouazzani@email.com", 32, "female", "123456789"},
		{"Adil", "Qadiri", "adil_q", "adil.qadiri@email.com", 29, "male", "123456789"},
		{"Kenza", "Rachidi", "kenza_r", "kenza.rachidi@email.com", 27, "female", "123456789"},
		{"Saad", "Sbihi", "saad_s", "saad.sbihi@email.com", 31, "male", "123456789"},
		{"Houda", "Tazi", "houda_t", "houda.tazi@email.com", 25, "female", "123456789"},
		{"Younes", "Wahbi", "younes_w", "younes.wahbi@email.com", 30, "male", "123456789"},
		{"Meryem", "Ziani", "meryem_z", "meryem.ziani@email.com", 28, "female", "123456789"},
	}
	

	for _, user := range users {
		if err := InsertUser(user); err != nil {
			log.Println(err)
		}
	}

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

// InsertUser inserts a user into the users table
func InsertUser(user RegistrationRequest) error {
	// Hash the password before inserting
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("could not hash password: %v", err)
	}

	query := `INSERT INTO users (first_name, last_name, nickname, email, age, gender, password) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = DB.Exec(query, user.FirstName, user.LastName, user.Nickname, user.Email, user.Age, user.Gender, hashedPassword)
	if err != nil {
		return fmt.Errorf("could not insert user: %v", err)
	}
	return nil
}