package route

import (
	"clean-arch/app/service"
	"clean-arch/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(app *fiber.App, db *mongo.Database) {

	RegisterFileRoutes(app, db)

	// Alumni Auth routes
	app.Post("/alumni/register", func(c *fiber.Ctx) error {
		return service.RegisterAlumniService(c, db)
	})

	app.Post("/alumni/login", func(c *fiber.Ctx) error {
		return service.AlumniLoginService(c, db)
	})

	app.Get("/alumni/profile", middleware.AlumniAuthRequired(), func(c *fiber.Ctx) error {
		return service.GetAlumniProfileService(c, db)
	})

	// Alumni routes
	app.Get("/alumni", func(c *fiber.Ctx) error {
		return service.GetAllAlumniService(c, db)
	})

	app.Get("/alumni/trash", middleware.UserAuthRequired(), middleware.UserOrAdminOnly(), func(c *fiber.Ctx) error {
		return service.GetTrashedAlumniService(c, db)
	})

	app.Get("/alumni/statistics", func(c *fiber.Ctx) error {
		return service.GetAlumniStatisticsService(c, db)
	})

	app.Get("/alumni/:id", func(c *fiber.Ctx) error {
		return service.GetAlumniByIDService(c, db)
	})

	app.Post("/alumni", func(c *fiber.Ctx) error {
		return service.CreateAlumniService(c, db)
	})

	app.Put("/alumni/:id", func(c *fiber.Ctx) error {
		return service.UpdateAlumniService(c, db)
	})

	app.Delete("/alumni/:id", func(c *fiber.Ctx) error {
		return service.DeleteAlumniService(c, db)
	})

	app.Post("/alumni/:id/soft-delete", middleware.UserAuthRequired(), middleware.UserOrAdminOnly(), func(c *fiber.Ctx) error {
		return service.SoftDeleteAlumniService(c, db)
	})

	app.Post("/alumni/:id/restore", middleware.UserAuthRequired(), middleware.UserOrAdminOnly(), func(c *fiber.Ctx) error {
		return service.RestoreAlumniService(c, db)
	})

	app.Delete("/alumni/:id/permanent", middleware.UserAuthRequired(), middleware.UserOrAdminOnly(), func(c *fiber.Ctx) error {
		return service.HardDeleteAlumniService(c, db)
	})

	app.Get("/cleanarch/alumni", func(c *fiber.Ctx) error {
		return service.GetAllAlumniWithPaginationService(c, db)
	})

	// Pekerjaan routes
	app.Get("/pekerjaan", func(c *fiber.Ctx) error {
		return service.GetAllPekerjaanService(c, db)
	})

	app.Get("/pekerjaan/:id", func(c *fiber.Ctx) error {
		return service.GetPekerjaanByIDService(c, db)
	})

	app.Get("/pekerjaan/alumni/:alumni_id", func(c *fiber.Ctx) error {
		return service.GetPekerjaanByAlumniIDService(c, db)
	})

	app.Post("/pekerjaan", middleware.AlumniAuthRequired(), func(c *fiber.Ctx) error {
		return service.CreatePekerjaanService(c, db)
	})

	app.Put("/pekerjaan/:id", middleware.AlumniAuthRequired(), func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanService(c, db)
	})

	app.Delete("/pekerjaan/:id", func(c *fiber.Ctx) error {
		return service.DeletePekerjaanService(c, db)
	})

	app.Get("/cleanarch/pekerjaan", func(c *fiber.Ctx) error {
		return service.GetAllPekerjaanWithPaginationService(c, db)
	})

	app.Delete("/pekerjaan/:id/soft", middleware.AlumniAuthRequired(), func(c *fiber.Ctx) error {
		return service.SoftDeletePekerjaanService(c, db)
	})

	// Original check route
	app.Post("/check/:key", func(c *fiber.Ctx) error {
		return service.CheckAlumniService(c, db)
	})

	// User Auth routes (for admin/system users)
	app.Post("/auth/login", func(c *fiber.Ctx) error {
		return service.LoginService(c, db)
	})

	app.Get("/auth/profile", middleware.UserAuthRequired(), func(c *fiber.Ctx) error {
		return service.GetProfileService(c, db)
	})
}
