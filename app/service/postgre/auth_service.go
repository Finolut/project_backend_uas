package service

import (
	"clean-arch/app/model/postgre"
	"clean-arch/app/repository/postgre"
	"clean-arch/utils/postgre"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func LoginService(c *fiber.Ctx, db *sql.DB) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Request body tidak valid",
		})
	}

	// Validasi input
	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Username dan password harus diisi",
		})
	}

	// Cari user di database
	user, passwordHash, err := repository.GetUserByUsernameOrEmail(db, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(401).JSON(fiber.Map{
				"error": "Username atau password salah",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Error database",
		})
	}

	// Check password
	if !utils.CheckPassword(req.Password, passwordHash) {
		return c.Status(401).JSON(fiber.Map{
			"error": "Username atau password salah",
		})
	}

	// Generate JWT token
	token, err := utils.GenerateToken(*user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal generate token",
		})
	}

	response := model.LoginResponse{
		User:  *user,
		Token: token,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data":    response,
	})
}

func AlumniLoginService(c *fiber.Ctx, db *sql.DB) error {
	var req model.AlumniLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Request body tidak valid",
		})
	}

	// Validasi input
	if req.NIM == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "NIM dan password harus diisi",
		})
	}

	// Cari alumni di database
	alumni, err := repository.GetAlumniByNIM(db, req.NIM)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(401).JSON(fiber.Map{
				"error": "NIM atau password salah",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Error database",
		})
	}

	// Check password
	if !utils.CheckPassword(req.Password, alumni.Password) {
		return c.Status(401).JSON(fiber.Map{
			"error": "NIM atau password salah",
		})
	}

	// Generate JWT token
	token, err := utils.GenerateAlumniToken(*alumni)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal generate token",
		})
	}

	// Remove password from response
	alumni.Password = ""

	response := model.AlumniLoginResponse{
		Alumni: *alumni,
		Token:  token,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data":    response,
	})
}

func GetAlumniProfileService(c *fiber.Ctx, db *sql.DB) error {
	alumniID := c.Locals("alumni_id").(int)

	// Get alumni with job history
	alumniWithJobs, err := repository.GetAlumniWithJobs(db, alumniID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil data profile",
		})
	}

	// Remove password from response
	alumniWithJobs.Alumni.Password = ""

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile berhasil diambil",
		"data":    alumniWithJobs,
	})
}

func RegisterAlumniService(c *fiber.Ctx, db *sql.DB) error {
	var req model.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Request body tidak valid",
		})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal hash password",
		})
	}

	// Create alumni
	alumni, err := repository.CreateAlumniWithAuth(db, req, hashedPassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal membuat akun alumni",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Akun alumni berhasil dibuat",
		"data":    alumni,
	})
}

func GetProfileService(c *fiber.Ctx, db *sql.DB) error {
	userID := c.Locals("user_id").(int)
	username := c.Locals("username").(string)
	role := c.Locals("role").(string)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile berhasil diambil",
		"data": fiber.Map{
			"user_id":  userID,
			"username": username,
			"role":     role,
		},
	})
}
