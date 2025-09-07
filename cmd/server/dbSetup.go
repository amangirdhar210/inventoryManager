package main

import (
	"database/sql"
	"fmt"
	"log"
)

func SetupDatabase(dbName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS products(
	"id" TEXT NOT NULL PRIMARY KEY,
	"name" TEXT,
	"price" REAL,
	"quantity" INTEGER
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("could not create products table %w", err)
	}

	log.Println("Database Initialized and Table created successfully.")
	return db, nil
}
