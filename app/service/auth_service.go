package service

import (
	"clean-arch/app/model"
	"clean-arch/app/repository"
	"clean-arch/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// LoginService godoc
// @Summary Login user admin/sistem
// @Description Melakukan login dengan username dan password untuk user admin atau sistem
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "Credentials login"
// @Success 200 {object} map[string]interface{} "Login berhasil dengan token"
// @Failure 400 {object} map[string]interface{} "Request body tidak valid atau field kosong"
// @Failure 401 {object} map[string]interface{} "Username atau password salah"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func LoginService(c *fiber.Ctx, db *mongo.Database) error {
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
		if err == mongo.ErrNoDocuments {
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

// AlumniLoginService godoc
// @Summary Login alumni
// @Description Melakukan login alumni berdasarkan NIM dan password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body model.AlumniLoginRequest true "NIM dan password alumni"
// @Success 200 {object} map[string]interface{} "Login berhasil dengan token"
// @Failure 400 {object} map[string]interface{} "Request body tidak valid atau field kosong"
// @Failure 401 {object} map[string]interface{} "NIM atau password salah"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alumni/login [post]
func AlumniLoginService(c *fiber.Ctx, db *mongo.Database) error {
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
		if err == mongo.ErrNoDocuments {
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

// GetAlumniProfileService godoc
// @Summary Dapatkan profile alumni
// @Description Mengambil profile lengkap alumni beserta riwayat pekerjaan
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Profile alumni dengan riwayat pekerjaan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security Bearer
// @Router /alumni/profile [get]
func GetAlumniProfileService(c *fiber.Ctx, db *mongo.Database) error {
	alumniID := c.Locals("alumni_id").(string)

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

// RegisterAlumniService godoc
// @Summary Register alumni baru
// @Description Membuat akun alumni baru dengan username/email dan password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body model.CreateAlumniRequest true "Data alumni untuk registrasi"
// @Success 200 {object} map[string]interface{} "Alumni berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "Request body tidak valid"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alumni/register [post]
func RegisterAlumniService(c *fiber.Ctx, db *mongo.Database) error {
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

// GetProfileService godoc
// @Summary Dapatkan profile user (admin)
// @Description Mengambil profile user admin yang sedang login
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Profile user admin"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security Bearer
// @Router /auth/profile [get]
func GetProfileService(c *fiber.Ctx, db *mongo.Database) error {
	userID := c.Locals("user_id").(string)
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
