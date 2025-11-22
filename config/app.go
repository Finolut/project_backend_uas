package config

import (
	"log"

	"clean-arch/app/service"
	"clean-arch/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewApp(db *mongo.Database) *fiber.App {
	app := fiber.New()
	
	app.Use(cors.New())
	app.Use(middleware.LoggerMiddleware)

	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	api := app.Group("/api")
	
	api.Get("/checkpoint", func(c *fiber.Ctx) error {
		return service.CheckpointService(c, db)
	})

	api.Post("/login", func(c *fiber.Ctx) error {
		return service.LoginService(c, db)
	})

	protected := api.Group("", middleware.AuthRequired())
	protected.Get("/profile", func(c *fiber.Ctx) error {
		return service.GetProfileService(c, db)
	})

	alumni := protected.Group("/alumni")
	alumni.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllAlumniService(c, db)
	})
	alumni.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetAlumniByIDService(c, db)
	})
	alumni.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreateAlumniService(c, db)
	})
	alumni.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdateAlumniService(c, db)
	})
	alumni.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.DeleteAlumniService(c, db)
	})

	alumni.Get("/trash", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.GetTrashedAlumniService(c, db)
	})
	alumni.Post("/:id/soft-delete", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.SoftDeleteAlumniService(c, db)
	})
	alumni.Post("/:id/restore", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.RestoreAlumniService(c, db)
	})
	alumni.Delete("/:id/permanent", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.HardDeleteAlumniService(c, db)
	})

	pekerjaan := protected.Group("/pekerjaan")
	pekerjaan.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllPekerjaanService(c, db)
	})
	pekerjaan.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetPekerjaanByIDService(c, db)
	})
	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.GetPekerjaanByAlumniIDService(c, db)
	})
	pekerjaan.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.CreatePekerjaanService(c, db)
	})
	pekerjaan.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanService(c, db)
	})
	pekerjaan.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.DeletePekerjaanService(c, db)
	})

	cleanarch := protected.Group("/cleanarch")
	log.Println("Registering /api/cleanarch/alumni route")
	cleanarch.Get("/alumni", func(c *fiber.Ctx) error {
		log.Println("Accessing /api/cleanarch/alumni endpoint")
		return service.GetAllAlumniWithPaginationService(c, db)
	})
	log.Println("Registering /api/cleanarch/pekerjaan route")
	cleanarch.Get("/pekerjaan", func(c *fiber.Ctx) error {
		log.Println("Accessing /api/cleanarch/pekerjaan endpoint")
		return service.GetAllPekerjaanWithPaginationService(c, db)
	})

	// Legacy route for compatibility
	app.Post("/check/:key", func(c *fiber.Ctx) error {
		return service.CheckAlumniService(c, db)
	})

	log.Println("All routes registered successfully:")
	log.Println("- POST /api/login")
	log.Println("- GET /api/profile (protected)")
	log.Println("- GET /api/alumni (protected)")
	log.Println("- GET /api/alumni/:id (protected)")
	log.Println("- POST /api/alumni (admin only)")
	log.Println("- PUT /api/alumni/:id (admin only)")
	log.Println("- DELETE /api/alumni/:id (admin only)")
	log.Println("- GET /api/alumni/trash (admin only)")
	log.Println("- POST /api/alumni/:id/soft-delete (admin only)")
	log.Println("- POST /api/alumni/:id/restore (admin only)")
	log.Println("- DELETE /api/alumni/:id/permanent (admin only)")
	log.Println("- GET /api/pekerjaan (protected)")
	log.Println("- GET /api/pekerjaan/:id (protected)")
	log.Println("- GET /api/pekerjaan/alumni/:alumni_id (admin only)")
	log.Println("- POST /api/pekerjaan (admin only)")
	log.Println("- PUT /api/pekerjaan/:id (admin only)")
	log.Println("- DELETE /api/pekerjaan/:id (admin only)")
	log.Println("- GET /api/cleanarch/alumni (protected, with pagination)")
	log.Println("- GET /api/cleanarch/pekerjaan (protected, with pagination)")
	log.Println("- POST /check/:key (legacy)")

	return app
}
