package service

import (
	"log"
	"os"
	"strconv"
	"strings"

	"clean-arch/app/model"
	"clean-arch/app/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// CheckAlumniService godoc
// @Summary Cek apakah mahasiswa adalah alumni
// @Description Mengecek keberadaan alumni berdasarkan NIM dengan API key
// @Tags Alumni
// @Accept mpform
// @Produce json
// @Param key path string true "API Key"
// @Param nim formData string true "NIM Mahasiswa"
// @Success 200 {object} map[string]interface{} "Alumni ditemukan atau tidak"
// @Failure 400 {object} map[string]interface{} "NIM wajib diisi"
// @Failure 401 {object} map[string]interface{} "Key tidak valid"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /check/{key} [post]
func CheckAlumniService(c *fiber.Ctx, db *mongo.Database) error {
	key := c.Params("key")
	if key != os.Getenv("API_KEY") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Key tidak valid",
			"success": false,
		})
	}

	nim := c.FormValue("nim")
	if nim == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "NIM wajib diisi",
			"success": false,
		})
	}

	alumni, err := repository.CheckAlumniByNim(db, nim)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message":  "Mahasiswa bukan alumni",
				"success":  true,
				"isAlumni": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal cek alumni karena " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Berhasil mendapatkan data alumni",
		"success":  true,
		"isAlumni": true,
		"alumni":   alumni,
	})
}

// GetAllAlumniService godoc
// @Summary Dapatkan semua alumni
// @Description Mengambil daftar semua alumni dari database
// @Tags Alumni
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Daftar alumni"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alumni [get]
func GetAllAlumniService(c *fiber.Ctx, db *mongo.Database) error {
	username := c.Locals("username").(string)
	log.Printf("User %s mengakses GET /api/alumni", username)

	alumni, err := repository.GetAllAlumni(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data alumni: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data alumni",
		"success": true,
		"data":    alumni,
	})
}

// GetAlumniByIDService godoc
// @Summary Dapatkan alumni berdasarkan ID
// @Description Mengambil data alumni spesifik berdasarkan ID MongoDB
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "Alumni ID (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{} "Data alumni"
// @Failure 404 {object} map[string]interface{} "Alumni tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alumni/{id} [get]
func GetAlumniByIDService(c *fiber.Ctx, db *mongo.Database) error {
	username := c.Locals("username").(string)
	id := c.Params("id")

	log.Printf("User %s mengakses GET /api/alumni/%s", username, id)

	alumni, err := repository.GetAlumniByID(db, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Alumni tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data alumni: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data alumni",
		"success": true,
		"data":    alumni,
	})
}

// CreateAlumniService godoc
// @Summary Buat alumni baru
// @Description Membuat data alumni baru di database
// @Tags Alumni
// @Accept json
// @Produce json
// @Param body body model.CreateAlumniRequest true "Data alumni baru"
// @Success 201 {object} map[string]interface{} "Alumni berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "Data tidak valid"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alumni [post]
func CreateAlumniService(c *fiber.Ctx, db *mongo.Database) error {
	username := c.Locals("username").(string)
	log.Printf("Admin %s menambah alumni baru", username)

	var req model.CreateAlumniRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format data tidak valid: " + err.Error(),
			"success": false,
		})
	}

	// Validasi input
	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "NIM, nama, jurusan, dan email wajib diisi",
			"success": false,
		})
	}

	alumni, err := repository.CreateAlumni(db, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menambah alumni: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Alumni berhasil ditambahkan",
		"success": true,
		"data":    alumni,
	})
}

// UpdateAlumniService godoc
// @Summary Update data alumni
// @Description Memperbarui data alumni berdasarkan ID
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "Alumni ID (MongoDB ObjectID)"
// @Param body body model.UpdateAlumniRequest true "Data alumni yang diupdate"
// @Success 200 {object} map[string]interface{} "Alumni berhasil diupdate"
// @Failure 400 {object} map[string]interface{} "Data tidak valid"
// @Failure 404 {object} map[string]interface{} "Alumni tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alumni/{id} [put]
func UpdateAlumniService(c *fiber.Ctx, db *mongo.Database) error {
	username := c.Locals("username").(string)
	id := c.Params("id")

	log.Printf("Admin %s mengupdate alumni ID %s", username, id)

	var req model.UpdateAlumniRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format data tidak valid: " + err.Error(),
			"success": false,
		})
	}

	// Validasi input
	if req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Nama, jurusan, dan email wajib diisi",
			"success": false,
		})
	}

	alumni, err := repository.UpdateAlumni(db, id, req)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Alumni tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengupdate alumni: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Alumni berhasil diupdate",
		"success": true,
		"data":    alumni,
	})
}

// DeleteAlumniService godoc
// @Summary Hapus alumni
// @Description Menghapus data alumni berdasarkan ID
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "Alumni ID (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{} "Alumni berhasil dihapus"
// @Failure 404 {object} map[string]interface{} "Alumni tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alumni/{id} [delete]
func DeleteAlumniService(c *fiber.Ctx, db *mongo.Database) error {
	username := c.Locals("username").(string)
	id := c.Params("id")

	log.Printf("Admin %s menghapus alumni ID %s", username, id)

	err := repository.DeleteAlumni(db, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Alumni tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus alumni: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Alumni berhasil dihapus",
		"success": true,
	})
}

// GetAlumniStatisticsService godoc
// @Summary Dapatkan statistik alumni
// @Description Mengambil statistik alumni berdasarkan berbagai kategori
// @Tags Alumni
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Statistik alumni"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alumni/statistics [get]
func GetAlumniStatisticsService(c *fiber.Ctx, db *mongo.Database) error {
	stats, err := repository.GetAlumniStatistics(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil statistik alumni: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil statistik alumni",
		"success": true,
		"data":    stats,
	})
}

// GetAllAlumniWithPaginationService godoc
// @Summary Dapatkan alumni dengan pagination
// @Description Mengambil daftar alumni dengan dukungan pagination, sorting, dan pencarian
// @Tags Alumni
// @Accept json
// @Produce json
// @Param page query int false "Halaman (default: 1)"
// @Param limit query int false "Limit data per halaman (default: 10, max: 100)"
// @Param sortBy query string false "Field untuk sorting (default: created_at)"
// @Param order query string false "Urutan sorting asc/desc (default: desc)"
// @Param search query string false "Pencarian berdasarkan nama atau email"
// @Success 200 {object} map[string]interface{} "Daftar alumni dengan metadata pagination"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /cleanarch/alumni [get]
func GetAllAlumniWithPaginationService(c *fiber.Ctx, db *mongo.Database) error {
	username := c.Locals("username").(string)
	log.Printf("User %s mengakses GET /api/alumni dengan pagination", username)

	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "created_at")
	order := c.Query("order", "desc")
	search := c.Query("search", "")

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Create pagination params
	params := model.PaginationParams{
		Page:   page,
		Limit:  limit,
		SortBy: sortBy,
		Order:  strings.ToLower(order),
		Search: search,
	}

	// Get data with pagination
	alumni, total, err := repository.GetAllAlumniWithPagination(db, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data alumni: " + err.Error(),
			"success": false,
		})
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	// Create response with pagination metadata
	response := model.AlumniResponse{
		Data: alumni,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  totalPages,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data alumni",
		"success": true,
		"data":    response.Data,
		"meta":    response.Meta,
	})
}

// GetTrashedAlumniService godoc
// @Summary Dapatkan alumni yang dihapus (soft delete)
// @Description Mengambil daftar alumni yang telah di-soft delete
// @Tags Alumni
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Daftar alumni yang dihapus"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security Bearer
// @Router /alumni/trash [get]
func GetTrashedAlumniService(c *fiber.Ctx, db *mongo.Database) error {
	// admin only via route middleware
	list, err := repository.GetTrashedAlumni(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil trash: " + err.Error(),
			"success": false,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data trash",
		"success": true,
		"data":    list,
	})
}

// SoftDeleteAlumniService godoc
// @Summary Soft delete alumni
// @Description Menghapus alumni secara soft (data masih tersimpan, hanya ditandai dihapus)
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "Alumni ID (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{} "Alumni berhasil dipindahkan ke trash"
// @Failure 404 {object} map[string]interface{} "Alumni tidak ditemukan atau sudah dihapus"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security Bearer
// @Router /alumni/{id}/soft-delete [post]
func SoftDeleteAlumniService(c *fiber.Ctx, db *mongo.Database) error {
	idStr := c.Params("id")

	userID, ok := c.Locals("user_id").(string)
	var deletedByID *string
	if ok && userID != "" {
		deletedByID = &userID
	}

	if err := repository.SoftDeleteAlumni(db, idStr, deletedByID); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Data alumni tidak ditemukan atau sudah dihapus",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal soft delete: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Alumni berhasil dipindahkan ke trash",
		"success": true,
	})
}

// RestoreAlumniService godoc
// @Summary Restore alumni dari trash
// @Description Mengembalikan alumni yang telah di-soft delete
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "Alumni ID (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{} "Alumni berhasil direstorasi"
// @Failure 404 {object} map[string]interface{} "Alumni tidak ditemukan atau belum di-trash"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security Bearer
// @Router /alumni/{id}/restore [post]
func RestoreAlumniService(c *fiber.Ctx, db *mongo.Database) error {
	idStr := c.Params("id")

	if err := repository.RestoreAlumni(db, idStr); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Data tidak ditemukan atau belum di-trash",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal restore: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil direstorasi dari trash",
		"success": true,
	})
}

// HardDeleteAlumniService godoc
// @Summary Hard delete alumni permanen
// @Description Menghapus alumni secara permanen dari database (hanya data alumni yang ada di trash)
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "Alumni ID (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{} "Data alumni dihapus permanen"
// @Failure 400 {object} map[string]interface{} "Hapus permanen hanya untuk data alumni yang ada di trash"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security Bearer
// @Router /alumni/{id}/permanent [delete]
func HardDeleteAlumniService(c *fiber.Ctx, db *mongo.Database) error {
	idStr := c.Params("id")

	if err := repository.HardDeleteAlumni(db, idStr); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Hapus permanen hanya untuk data alumni yang ada di trash",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal hapus permanen: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data alumni dihapus permanen",
		"success": true,
	})
}
