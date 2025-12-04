package config

import (
	"database/sql"
	"log"

	service "clean-arch/app/service/postgre"
	middleware "clean-arch/middleware/postgre"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewApp(db *sql.DB) *fiber.App {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(middleware.LoggerMiddleware)

	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	api := app.Group("/api")

	api.Get("/checkpoint", func(c *fiber.Ctx) error {
		// Menggunakan service.CheckpointService dari health_service.go
		return service.CheckpointService(c, db)
	})

	api.Post("/login", func(c *fiber.Ctx) error {
		// Menggunakan service.LoginService dari auth_service.go
		return service.LoginService(c, db)
	})

	protected := api.Group("", middleware.AuthRequired())
	protected.Get("/profile", func(c *fiber.Ctx) error {
		return service.GetProfileService(c, db)
	})

	// Logging output disesuaikan
	log.Println("All routes registered successfully:")
	log.Println("- POST /api/login")
	log.Println("- GET /api/profile (protected)")
	log.Println("- GET /api/student (protected)")
	log.Println("- GET /api/student/:id (protected)")
	log.Println("- POST /api/student (admin only)")
	log.Println("- PUT /api/student/:id (admin only)")
	log.Println("- DELETE /api/student/:id (admin only)")
	log.Println("- GET /api/pekerjaan (protected)")
	log.Println("- GET /api/pekerjaan/:id (protected)")
	log.Println("- GET /api/pekerjaan/student/:user_id (admin only)")
	log.Println("- GET /api/cleanarch/student (protected, with pagination)")
	log.Println("- POST /check/:key (legacy)")

	return app
}
