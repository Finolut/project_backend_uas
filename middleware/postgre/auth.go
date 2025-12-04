package middleware

import (
	utils "clean-arch/utils/postgre"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Middleware untuk memerlukan login
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token akses diperlukan",
			})
		}

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
		c.Locals("role", claims.RoleName) // Menggunakan RoleName: "Admin", "Mahasiswa", dll.
		c.Locals("role_id", claims.RoleID)

		return c.Next()
	}
}

// Mengganti AlumniAuthRequired
func StudentAuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token akses diperlukan",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Format token tidak valid",
			})
		}

		// Menggunakan ValidateStudentToken
		claims, err := utils.ValidateStudentToken(tokenParts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak valid atau expired",
			})
		}

		// Simpan informasi student di context
		c.Locals("user_id", claims.UserID) // ID dari tabel users
		c.Locals("nim", claims.NIM)
		c.Locals("nama", claims.Nama)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.RoleName) // Role Name: "Mahasiswa"
		c.Locals("is_student", true)

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
		// Memeriksa nama role dari SRS: "Admin"
		if roleStr != "Admin" {
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
				"error": "Akses ditolak. Hanya mahasiswa atau admin yang diizinkan",
			})
		}
		roleStr := role.(string)
		// Memeriksa nama role dari SRS: "Mahasiswa" atau "Admin"
		if roleStr != "Mahasiswa" && roleStr != "Admin" {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses ditolak. Hanya mahasiswa atau admin yang diizinkan",
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
		// Mengasumsikan "user" lama adalah "Mahasiswa"
		if roleStr != "Mahasiswa" && roleStr != "Admin" && roleStr != "Dosen Wali" {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses ditolak. Hanya user/mahasiswa, dosen wali, atau admin yang diizinkan",
			})
		}
		return c.Next()
	}
}

// Middleware untuk user biasa (bukan alumni) - Menggunakan AuthRequired utama
func UserAuthRequired() fiber.Handler {
	return AuthRequired()
}
