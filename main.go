package main

import (
	"log"
	"os"

	// Import library eksternal
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	// Import Config
	mongoConfig "clean-arch/config/mongo"
	postgreConfig "clean-arch/config/postgre"

	// Import Database Connection
	mongoDB "clean-arch/database/mongo"
	postgreDB "clean-arch/database/postgre"

	// Import Docs (Swagger)
	_ "clean-arch/docs"

	// Import Routes
	mongoRoute "clean-arch/route/mongo"
	postgreRoute "clean-arch/route/postgre"
)

// @title Auth API v1
// @version 1.0
// @description API untuk autentikasi dan manajemen user menggunakan Clean Architecture (Support MongoDB & PostgreSQL)
// @host localhost:3000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

func main() {
	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}

	// Ambil konfigurasi port dan driver database
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Tentukan driver: "mongo" atau "postgres" (default ke mongo jika kosong)
	dbDriver := os.Getenv("DB_DRIVER")

	var app *fiber.App

	// 2. Inisialisasi berdasarkan Driver Database
	if dbDriver == "postgres" {
		log.Println("🐘 Starting application with PostgreSQL...")

		// a. Koneksi ke PostgreSQL
		db := postgreDB.ConnectDB()
		// defer db.Close() // Opsional: tergantung lifecycle aplikasi

		// b. Setup App (Middleware, Static files, dll khusus Postgre config)
		app = postgreConfig.NewApp(db)

		// c. Setup Swagger Route
		app.Get("/swagger/*", fiberSwagger.WrapHandler)

		// d. Register Routes khusus PostgreSQL
		postgreRoute.RegisterRoutes(app, db)

	} else {
		// Default: MongoDB
		log.Println("🍃 Starting application with MongoDB...")

		// a. Koneksi ke MongoDB
		client := mongoDB.ConnectDB()
		db := mongoDB.GetDatabase(client)
		defer mongoDB.DisconnectDB(client)

		// b. Setup App (Middleware, Static files, dll khusus Mongo config)
		app = mongoConfig.NewApp(db)

		// c. Setup Swagger Route
		app.Get("/swagger/*", fiberSwagger.WrapHandler)

		// d. Register Routes khusus MongoDB
		mongoRoute.RegisterRoutes(app, db)
	}

	// 3. Jalankan Server
	log.Printf("🚀 Server running on port %s using %s driver", port, dbDriver)
	log.Println("📚 Swagger docs available at http://localhost:" + port + "/swagger/index.html")
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
