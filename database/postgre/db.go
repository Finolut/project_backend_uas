package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Successfully connected to database")
	return db
}
