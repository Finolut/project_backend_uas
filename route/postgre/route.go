package route

import (
	"database/sql"

	service "clean-arch/app/service/postgre"
	middleware "clean-arch/middleware/postgre"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, db *sql.DB) {

	// --- API v1 Auth Routes ---
	authV1 := app.Group("/api/v1/auth")

	// Public auth endpoints
	authV1.Post("/login", func(c *fiber.Ctx) error {
		return service.LoginService(c, db)
	})

	authV1.Post("/refresh", func(c *fiber.Ctx) error {
		return service.RefreshTokenService(c, db)
	})

	// Protected auth endpoints
	authV1.Post("/logout", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		return service.LogoutService(c, db)
	})

	authV1.Get("/profile", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		return service.GetProfileService(c, db)
	})
}
