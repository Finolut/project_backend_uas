package middleware

import (
	"clean-arch/utils/mongo"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Middleware untuk memerlukan login
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token akses diperlukan",
			})
		}

		// Extract token dari "Bearer TOKEN"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Format token tidak valid",
			})
		}

		// Validasi token
		claims, err := utils.ValidateToken(tokenParts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak valid atau expired",
			})
		}

		// Simpan informasi user di context
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func AlumniAuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token akses diperlukan",
			})
		}

		// Extract token dari "Bearer TOKEN"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Format token tidak valid",
			})
		}

		// Validasi alumni token
		claims, err := utils.ValidateAlumniToken(tokenParts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak valid atau expired",
			})
		}

		// Simpan informasi alumni di context
		c.Locals("alumni_id", claims.AlumniID)
		c.Locals("nim", claims.NIM)
		c.Locals("nama", claims.Nama)
		c.Locals("email", claims.Email)
		c.Locals("role", "alumni")
		c.Locals("is_alumni", true)

		return c.Next()
	}
}

// Middleware untuk memerlukan role admin
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses ditolak. Hanya admin yang diizinkan",
			})
		}
		roleStr := role.(string)
		if roleStr != "admin" {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses ditolak. Hanya admin yang diizinkan",
			})
		}
		return c.Next()
	}
}

func AlumniOrAdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses ditolak. Hanya alumni atau admin yang diizinkan",
			})
		}
		roleStr := role.(string)
		if roleStr != "alumni" && roleStr != "admin" {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses ditolak. Hanya alumni atau admin yang diizinkan",
			})
		}
		return c.Next()
	}
}

func UserOrAdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses ditolak. Hanya user atau admin yang diizinkan",
			})
		}
		roleStr := role.(string)
		if roleStr != "user" && roleStr != "admin" {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses ditolak. Hanya user atau admin yang diizinkan",
			})
		}
		return c.Next()
	}
}

// Middleware untuk user biasa (bukan alumni)
func UserAuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token akses diperlukan",
			})
		}

		// Extract token dari "Bearer TOKEN"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Format token tidak valid",
			})
		}

		// Validasi token
		claims, err := utils.ValidateToken(tokenParts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak valid atau expired",
			})
		}

		// Simpan informasi user di context
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)
		c.Locals("is_user", true)

		return c.Next()
	}
}
