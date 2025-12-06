package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
)

// Use this only if you need CSRF protection for cookie/form based flows.
// For JWT-in-Authorization header APIs it's typically not required.
func CSRF() fiber.Handler {
	return csrf.New(csrf.Config{
		KeyLookup:      "header:X-CSRF-Token", // or "form:_csrf"
		CookieName:     "csrf_token",
		CookieSameSite: "Lax",
		Expiration:     24 * time.Hour,
	})
}
