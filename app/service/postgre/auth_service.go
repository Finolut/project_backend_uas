package service

import (
	model "clean-arch/app/model/postgre"
	repository "clean-arch/app/repository/postgre"
	utils "clean-arch/utils/postgre"
	"database/sql"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func LoginService(c *fiber.Ctx, db *sql.DB) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Request body tidak valid",
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Username dan password harus diisi",
		})
	}

	// Mendapatkan User, Password Hash, dan Role Name DARI TABEL USERS (General Login)
	user, passwordHash, roleName, err := repository.GetUserByUsernameOrEmail(db, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(401).JSON(fiber.Map{
				"error": "Username atau password salah",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Error database: " + err.Error(),
		})
	}

	if !user.IsActive {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akun tidak aktif",
		})
	}

	if !utils.CheckPassword(req.Password, passwordHash) {
		return c.Status(401).JSON(fiber.Map{
			"error": "Username atau password salah",
		})
	}

	token, err := utils.GenerateToken(*user, roleName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal generate token",
		})
	}

	// Prepare response user model (excluding hash)
	user.PasswordHash = ""

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

// StudentLoginService (Mengganti AlumniLoginService)
func StudentLoginService(c *fiber.Ctx, db *sql.DB) error {
	// Menggunakan LoginRequest yang sama untuk menerima Username/Email dan Password
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Request body tidak valid",
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Username dan password harus diisi",
		})
	}

	// 1. Mendapatkan User, Password Hash, dan Role Name DARI TABEL USERS
	user, passwordHash, roleName, err := repository.GetUserByUsernameOrEmail(db, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(401).JSON(fiber.Map{
				"error": "Username atau password salah",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Error database: " + err.Error(),
		})
	}

	// Verifikasi apakah user adalah mahasiswa (role harus "Mahasiswa")
	if roleName != "Mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak. Akun bukan akun Mahasiswa.",
		})
	}

	if !user.IsActive {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akun tidak aktif",
		})
	}

	if !utils.CheckPassword(req.Password, passwordHash) {
		return c.Status(401).JSON(fiber.Map{
			"error": "Username atau password salah",
		})
	}

	// 2. Ambil detail Student (yang berisi NIM) menggunakan user.ID
	student, err := repository.GetStudentByUserID(db, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Kasus ini seharusnya jarang terjadi jika user.Role adalah "Mahasiswa",
			// tapi ini untuk keamanan data integrity.
			return c.Status(500).JSON(fiber.Map{
				"error": "Data mahasiswa tidak lengkap. Hubungi Admin.",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Error mengambil data mahasiswa: " + err.Error(),
		})
	}

	// Menggunakan GenerateStudentToken
	token, err := utils.GenerateStudentToken(*user, *student, roleName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal generate token",
		})
	}

	// Buat response model StudentView
	studentView := model.StudentView{
		ID:        user.ID,
		NIM:       student.StudentID,
		Nama:      user.FullName,
		Email:     user.Email,
		Role:      roleName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	response := model.StudentLoginResponse{
		Student: studentView,
		Token:   token,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data":    response,
	})
}

// GetStudentProfileService (Mengganti GetAlumniProfileService)
func GetStudentProfileService(c *fiber.Ctx, db *sql.DB) error {
	userID := c.Locals("user_id").(string)

	user, err := repository.GetUserByID(db, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data user: " + err.Error()})
	}

	student, err := repository.GetStudentByUserID(db, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data mahasiswa: " + err.Error()})
	}

	user.PasswordHash = ""

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile berhasil diambil",
		"data": fiber.Map{
			"user":            *user,
			"student_profile": student,
		},
	})
}

// RegisterStudentService (Mengganti RegisterAlumniService)
func RegisterStudentService(c *fiber.Ctx, db *sql.DB) error {
	var req model.RegisterStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Request body tidak valid",
		})
	}

	// Validasi input minimal
	if req.StudentID == "" || req.FullName == "" || req.Email == "" || req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "Semua field wajib diisi",
			"success": false,
		})
	}

	// Find Student Role ID
	role, err := repository.GetRoleByName(db, "Mahasiswa")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Role Mahasiswa tidak ditemukan"})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal hash password"})
	}

	// Create User and Student
	user, student, err := repository.CreateUserAndStudent(db, req, hashedPassword, role.ID)
	if err != nil {
		if strings.Contains(err.Error(), "username atau email sudah terdaftar") || strings.Contains(err.Error(), "NIM sudah terdaftar") {
			return c.Status(409).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat akun mahasiswa: " + err.Error()})
	}

	user.PasswordHash = ""

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Akun mahasiswa berhasil dibuat",
		"data": fiber.Map{
			"user":            *user,
			"student_profile": *student,
		},
	})
}

func GetProfileService(c *fiber.Ctx, db *sql.DB) error {
	userID := c.Locals("user_id").(string)

	user, err := repository.GetUserByID(db, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil data user: " + err.Error(),
		})
	}

	user.PasswordHash = ""

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile berhasil diambil",
		"data":    *user,
	})
}
