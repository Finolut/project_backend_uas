package service

import (
	model "clean-arch/app/model/postgre"
	repository "clean-arch/app/repository/postgre"
	utils "clean-arch/utils/postgre"
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
)

// LoginService godoc
// @Summary User Login
// @Description Login untuk user (Admin/Sistem) dengan username dan password
// @Tags Auth (v1)
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "Username dan password"
// @Success 200 {object} map[string]interface{} "Login berhasil dengan access token dan refresh token"
// @Failure 400 {object} map[string]interface{} "Request body tidak valid"
// @Failure 401 {object} map[string]interface{} "Username atau password salah"
// @Failure 403 {object} map[string]interface{} "Akun tidak aktif"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/login [post]
func LoginService(c *fiber.Ctx, db *sql.DB) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Request body tidak valid",
			"success": false,
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Username dan password harus diisi",
			"success": false,
		})
	}

	user, passwordHash, roleName, err := repository.GetUserByUsernameOrEmail(db, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(401).JSON(fiber.Map{
				"error":   "Username atau password salah",
				"success": false,
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error":   "Database error",
			"success": false,
		})
	}

	if !user.IsActive {
		return c.Status(403).JSON(fiber.Map{
			"error":   "Akun tidak aktif",
			"success": false,
		})
	}

	if !utils.CheckPassword(req.Password, passwordHash) {
		return c.Status(401).JSON(fiber.Map{
			"error":   "Username atau password salah",
			"success": false,
		})
	}

	accessToken, err := utils.GenerateToken(*user, roleName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Gagal membuat access token",
			"success": false,
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(*user, roleName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Gagal membuat refresh token",
			"success": false,
		})
	}

	user.PasswordHash = ""

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data": fiber.Map{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    3600,
		},
	})
}

// RefreshTokenService godoc
// @Summary Refresh Access Token
// @Description Menggunakan refresh token untuk mendapatkan access token baru
// @Tags Auth (v1)
// @Accept json
// @Produce json
// @Param body body model.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]interface{} "Token refresh berhasil"
// @Failure 400 {object} map[string]interface{} "Request body tidak valid"
// @Failure 401 {object} map[string]interface{} "Refresh token tidak valid atau expired"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/refresh [post]
func RefreshTokenService(c *fiber.Ctx, db *sql.DB) error {
	var req model.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Request body tidak valid",
			"success": false,
		})
	}

	if req.RefreshToken == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Refresh token harus diisi",
			"success": false,
		})
	}

	claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error":   "Refresh token tidak valid atau expired",
			"success": false,
		})
	}

	user, err := repository.GetUserByID(db, claims.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(401).JSON(fiber.Map{
				"error":   "User tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error":   "Database error",
			"success": false,
		})
	}

	newAccessToken, err := utils.GenerateToken(*user, claims.RoleName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Gagal membuat access token",
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Token berhasil di-refresh",
		"data": fiber.Map{
			"access_token": newAccessToken,
			"token_type":   "Bearer",
			"expires_in":   3600,
		},
	})
}

// LogoutService godoc
// @Summary User Logout
// @Description Logout user dan invalidate token (opsional jika menggunakan blacklist)
// @Tags Auth (v1)
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Logout berhasil"
// @Failure 400 {object} map[string]interface{} "Request tidak valid"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/logout [post]
// @Security Bearer
func LogoutService(c *fiber.Ctx, db *sql.DB) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "User ID tidak ditemukan di token",
			"success": false,
		})
	}

	// Implementasi logout bisa menggunakan:
	// 1. Blacklist token (store di Redis atau database)
	// 2. Invalidate refresh token
	// 3. Update last_logout di database
	// Untuk sekarang, hanya return success message

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logout berhasil",
	})
}

// GetProfileService godoc
// @Summary Get User Profile
// @Description Mengambil profile user yang sedang login berdasarkan token
// @Tags Auth (v1)
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Profile user"
// @Failure 401 {object} map[string]interface{} "Token tidak valid"
// @Failure 404 {object} map[string]interface{} "User tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/profile [get]
// @Security Bearer
func GetProfileService(c *fiber.Ctx, db *sql.DB) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(401).JSON(fiber.Map{
			"error":   "Token tidak valid",
			"success": false,
		})
	}

	user, err := repository.GetUserByID(db, userID.(string))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(404).JSON(fiber.Map{
				"error":   "User tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error":   "Database error",
			"success": false,
		})
	}

	user.PasswordHash = ""

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile berhasil diambil",
		"data":    user,
	})
}
