package route

import (
	"clean-arch/app/service/mongo"
	"clean-arch/middleware/mongo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterFileRoutes registers all file upload routes
func RegisterFileRoutes(app *fiber.App, db *mongo.Database) {
	files := app.Group("/api/files")

	// POST /api/files/upload-photo
	// Requires: user token (admin or regular user)
	// Body: form-data with file (max 1MB) and user_id
	files.Post("/upload-photo", middleware.FileAuthRequired(), func(c *fiber.Ctx) error {
		return service.UploadPhotoService(c, db)
	})

	// POST /api/files/upload-certificate
	// Requires: user token (admin or regular user)
	// Body: form-data with file (max 2MB PDF) and user_id
	files.Post("/upload-certificate", middleware.FileAuthRequired(), func(c *fiber.Ctx) error {
		return service.UploadCertificateService(c, db)
	})

	// GET /api/files?user_id=xxx&category=photo|certificate
	// Requires: user token (admin or regular user)
	files.Get("/", middleware.FileAuthRequired(), func(c *fiber.Ctx) error {
		return service.GetFilesService(c, db)
	})

	// DELETE /api/files/:id
	// Requires: user token (admin or regular user)
	files.Delete("/:id", middleware.FileAuthRequired(), func(c *fiber.Ctx) error {
		return service.DeleteFileService(c, db)
	})
}
