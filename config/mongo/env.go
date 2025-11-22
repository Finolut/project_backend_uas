package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Coba load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env file not found, using system environment variables")
	}

	// Default values jika belum ada
	if os.Getenv("APP_PORT") == "" {
		os.Setenv("APP_PORT", "3000")
	}
	if os.Getenv("API_KEY") == "" {
		os.Setenv("API_KEY", "12345")
	}
	// MongoDB configuration
	if os.Getenv("MONGODB_URI") == "" {
		os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	}
	if os.Getenv("MONGODB_DB") == "" {
		os.Setenv("MONGODB_DB", "alumni_db")
	}
	// Keep for backward compatibility
	if os.Getenv("DB_DSN") == "" {
		os.Setenv("DB_DSN", "postgres://postgres:admin@localhost:5432/alumni_db?sslmode=disable")
	}
}
