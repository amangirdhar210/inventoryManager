package main

import (
	"database/sql"
	"log"

	"github.com/amangirdhar210/inventory-manager/internal/core/domain"
	"github.com/google/uuid"
)

func SetupDatabase(dbName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	createProductsTableSQL := `
    CREATE TABLE IF NOT EXISTS products(
        "id" TEXT NOT NULL PRIMARY KEY,
        "name" TEXT,
        "price" REAL,
        "quantity" INTEGER
    );`
	if _, err := db.Exec(createProductsTableSQL); err != nil {
		return nil, err
	}

	createManagersTableSQL := `
    CREATE TABLE IF NOT EXISTS managers(
        "id" TEXT NOT NULL PRIMARY KEY,
        "email" TEXT UNIQUE,
        "password" TEXT
    );`
	if _, err := db.Exec(createManagersTableSQL); err != nil {
		return nil, err
	}

	seedAdmin(db)

	log.Println("Database Initialized and Tables created successfully.")
	return db, nil
}

func seedAdmin(db *sql.DB) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM managers WHERE email = ?", "admin@example.com")
	if err := row.Scan(&count); err != nil {
		log.Printf("Could not check for admin user: %v", err)
		return
	}

	if count == 0 {
		admin := &domain.Manager{
			Id:       uuid.NewString(),
			Email:    "admin@example.com",
			Password: "password123",
		}
		if err := admin.HashPassword(); err != nil {
			log.Printf("Could not hash admin password: %v", err)
			return
		}

		stmt, err := db.Prepare("INSERT INTO managers(id, email, password) VALUES(?,?,?)")
		if err != nil {
			log.Printf("Could not prepare admin insert statement: %v", err)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(admin.Id, admin.Email, admin.Password)
		if err != nil {
			log.Printf("Could not seed admin user: %v", err)
			return
		}
		log.Println("Admin user seeded successfully.")
		log.Println("Email: admin@example.com")
		log.Println("Password: password123")
	}
}
