package main

import (
	"log"
	"os"

	"clean-arch/config/mongo"
	"clean-arch/database"
	_ "clean-arch/docs"
	"clean-arch/route"

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Alumni Management API
// @version 1.0
// @description API untuk mengelola data alumni dengan MongoDB menggunakan Clean Architecture
// @host localhost:3000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

func main() {
	config.LoadEnv()
	client := database.ConnectDB()
	db := database.GetDatabase(client)
	defer database.DisconnectDB(client)

	app := config.NewApp(db)

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	route.RegisterRoutes(app, db)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
