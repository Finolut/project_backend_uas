package service

import (
	"log"
	"strconv"
	"strings"

	"clean-arch/app/model"
	"clean-arch/app/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllPekerjaanService godoc
// @Summary Dapatkan semua riwayat pekerjaan
// @Description Mengambil daftar semua riwayat pekerjaan alumni
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Daftar riwayat pekerjaan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /pekerjaan [get]
func GetAllPekerjaanService(c *fiber.Ctx, db *mongo.Database) error {
	username := c.Locals("username").(string)
	log.Printf("User %s mengakses GET /api/pekerjaan", username)

	pekerjaan, err := repository.GetAllPekerjaan(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data pekerjaan",
		"success": true,
		"data":    pekerjaan,
	})
}

// GetPekerjaanByIDService godoc
// @Summary Dapatkan riwayat pekerjaan berdasarkan ID
// @Description Mengambil detail riwayat pekerjaan berdasarkan ID
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "Pekerjaan ID (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{} "Detail riwayat pekerjaan"
// @Failure 404 {object} map[string]interface{} "Pekerjaan tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /pekerjaan/{id} [get]
func GetPekerjaanByIDService(c *fiber.Ctx, db *mongo.Database) error {
	username := c.Locals("username").(string)
	id := c.Params("id")

	log.Printf("User %s mengakses GET /api/pekerjaan/%s", username, id)

	pekerjaan, err := repository.GetPekerjaanByID(db, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Pekerjaan tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data pekerjaan",
		"success": true,
		"data":    pekerjaan,
	})
}

// GetPekerjaanByAlumniIDService godoc
// @Summary Dapatkan riwayat pekerjaan alumni
// @Description Mengambil semua riwayat pekerjaan berdasarkan alumni ID
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param alumni_id path string true "Alumni ID (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{} "Daftar riwayat pekerjaan alumni"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /pekerjaan/alumni/{alumni_id} [get]
func GetPekerjaanByAlumniIDService(c *fiber.Ctx, db *mongo.Database) error {
	username := c.Locals("username").(string)
	alumniID := c.Params("alumni_id")

	log.Printf("Admin %s mengakses GET /api/pekerjaan/alumni/%s", username, alumniID)

	pekerjaan, err := repository.GetPekerjaanByAlumniID(db, alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data pekerjaan alumni",
		"success": true,
		"data":    pekerjaan,
	})
}

// CreatePekerjaanService godoc
// @Summary Buat riwayat pekerjaan baru
// @Description Membuat riwayat pekerjaan baru untuk alumni
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param body body model.CreatePekerjaanRequest true "Data riwayat pekerjaan"
// @Success 201 {object} map[string]interface{} "Riwayat pekerjaan berhasil dibuat"
// @Failure 400 {object} map[string]interface{} "Data tidak valid"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security Bearer
// @Router /pekerjaan [post]
func CreatePekerjaanService(c *fiber.Ctx, db *mongo.Database) error {
	alumniID := c.Locals("alumni_id").(string)
	nama := c.Locals("nama").(string)
	log.Printf("Alumni %s (ID: %s) menambah pekerjaan baru", nama, alumniID)

	var req model.CreatePekerjaanRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format data tidak valid: " + err.Error(),
			"success": false,
		})
	}

	// Validasi input
	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Nama perusahaan dan posisi jabatan wajib diisi",
			"success": false,
		})
	}

	pekerjaan, err := repository.CreatePekerjaan(db, req, alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menambah pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pekerjaan berhasil ditambahkan",
		"success": true,
		"data":    pekerjaan,
	})
}

// UpdatePekerjaanService godoc
// @Summary Update riwayat pekerjaan
// @Description Memperbarui riwayat pekerjaan (hanya untuk pemilik atau admin)
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "Pekerjaan ID (MongoDB ObjectID)"
// @Param body body model.UpdatePekerjaanRequest true "Data riwayat pekerjaan yang diupdate"
// @Success 200 {object} map[string]interface{} "Riwayat pekerjaan berhasil diupdate"
// @Failure 400 {object} map[string]interface{} "Data tidak valid"
// @Failure 403 {object} map[string]interface{} "Anda tidak punya akses"
// @Failure 404 {object} map[string]interface{} "Pekerjaan tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security Bearer
// @Router /pekerjaan/{id} [put]
func UpdatePekerjaanService(c *fiber.Ctx, db *mongo.Database) error {
	alumniID := c.Locals("alumni_id").(string)
	nama := c.Locals("nama").(string)
	id := c.Params("id")

	log.Printf("Alumni %s (ID: %s) mengupdate pekerjaan ID %s", nama, alumniID, id)

	ownerAlumniID, err := repository.GetAlumniIDByPekerjaanID(db, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Pekerjaan tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	if ownerAlumniID != alumniID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Anda hanya dapat mengupdate riwayat pekerjaan milik Anda sendiri",
			"success": false,
		})
	}

	var req model.UpdatePekerjaanRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format data tidak valid: " + err.Error(),
			"success": false,
		})
	}

	// Validasi input
	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Nama perusahaan dan posisi jabatan wajib diisi",
			"success": false,
		})
	}

	pekerjaan, err := repository.UpdatePekerjaan(db, id, req)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Pekerjaan tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengupdate pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan berhasil diupdate",
		"success": true,
		"data":    pekerjaan,
	})
}

// DeletePekerjaanService godoc
// @Summary Hapus riwayat pekerjaan
// @Description Menghapus riwayat pekerjaan secara permanent
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "Pekerjaan ID (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{} "Riwayat pekerjaan berhasil dihapus"
// @Failure 404 {object} map[string]interface{} "Pekerjaan tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /pekerjaan/{id} [delete]
func DeletePekerjaanService(c *fiber.Ctx, db *mongo.Database) error {
	username := c.Locals("username").(string)
	id := c.Params("id")

	log.Printf("Admin %s menghapus pekerjaan ID %s", username, id)

	err := repository.DeletePekerjaan(db, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Pekerjaan tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan berhasil dihapus",
		"success": true,
	})
}

// GetAllPekerjaanWithPaginationService godoc
// @Summary Dapatkan riwayat pekerjaan dengan pagination
// @Description Mengambil daftar riwayat pekerjaan dengan dukungan pagination, sorting, dan pencarian
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param page query int false "Halaman (default: 1)"
// @Param limit query int false "Limit data per halaman (default: 10, max: 100)"
// @Param sortBy query string false "Field untuk sorting (default: created_at)"
// @Param order query string false "Urutan sorting asc/desc (default: desc)"
// @Param search query string false "Pencarian berdasarkan nama perusahaan atau jabatan"
// @Success 200 {object} map[string]interface{} "Daftar riwayat pekerjaan dengan metadata pagination"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /cleanarch/pekerjaan [get]
func GetAllPekerjaanWithPaginationService(c *fiber.Ctx, db *mongo.Database) error {
	username := c.Locals("username").(string)
	log.Printf("User %s mengakses GET /api/pekerjaan dengan pagination", username)

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
	pekerjaan, total, err := repository.GetAllPekerjaanWithPagination(db, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	// Create response with pagination metadata
	response := model.PekerjaanResponse{
		Data: pekerjaan,
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
		"message": "Berhasil mengambil data pekerjaan",
		"success": true,
		"data":    response.Data,
		"meta":    response.Meta,
	})
}

// SoftDeletePekerjaanService godoc
// @Summary Soft delete riwayat pekerjaan
// @Description Menghapus riwayat pekerjaan secara soft delete (data masih tersimpan)
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "Pekerjaan ID (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{} "Riwayat pekerjaan berhasil dihapus"
// @Failure 403 {object} map[string]interface{} "Anda tidak punya akses"
// @Failure 404 {object} map[string]interface{} "Pekerjaan tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security Bearer
// @Router /pekerjaan/{id}/soft [delete]
func SoftDeletePekerjaanService(c *fiber.Ctx, db *mongo.Database) error {
	alumniID, isAlumni := c.Locals("alumni_id").(string)
	role := c.Locals("role").(string)
	username := c.Locals("username").(string)
	id := c.Params("id")

	// Get the owner of this job record
	ownerAlumniID, err := repository.GetAlumniIDByPekerjaanID(db, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Pekerjaan tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	if isAlumni && role == "alumni" {
		// Alumni trying to delete - can only delete their own
		if ownerAlumniID != alumniID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Maaf tidak bisa mengubah atau mengedit data selain diri sendiri",
				"success": false,
			})
		}
		log.Printf("Alumni ID %s menghapus pekerjaan ID %s miliknya sendiri", alumniID, id)
	} else if role == "admin" {
		// Admin can delete any
		log.Printf("Admin %s menghapus pekerjaan ID %s", username, id)
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Anda tidak memiliki akses untuk menghapus pekerjaan",
			"success": false,
		})
	}

	var deletedBy string
	if role == "admin" {
		userID, ok := c.Locals("user_id").(string)
		if !ok {
			userID = ""
		}
		deletedBy = userID
	} else {
		deletedBy = alumniID
	}

	err = repository.SoftDeletePekerjaan(db, id, deletedBy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Pekerjaan tidak ditemukan atau sudah dihapus",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan berhasil dihapus",
		"success": true,
	})
}

// RestorePekerjaanService godoc
// @Summary Restore riwayat pekerjaan dari soft delete
// @Description Mengembalikan riwayat pekerjaan yang telah di-soft delete
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "Pekerjaan ID (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{} "Riwayat pekerjaan berhasil direstorasi"
// @Failure 403 {object} map[string]interface{} "Anda tidak punya akses"
// @Failure 404 {object} map[string]interface{} "Pekerjaan tidak ditemukan atau tidak di-trash"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security Bearer
// @Router /pekerjaan/{id}/restore [post]
func RestorePekerjaanService(c *fiber.Ctx, db *mongo.Database) error {
	alumniID, isAlumni := c.Locals("alumni_id").(string)
	role := c.Locals("role").(string)
	username := c.Locals("username").(string)
	id := c.Params("id")

	// Get the owner of this job record
	ownerAlumniID, err := repository.GetAlumniIDByPekerjaanID(db, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Pekerjaan tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	if isAlumni && role == "alumni" {
		// Alumni trying to restore - can only restore their own
		if ownerAlumniID != alumniID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Maaf tidak bisa mengubah atau mengedit data selain diri sendiri",
				"success": false,
			})
		}
		log.Printf("Alumni ID %s merestorasi pekerjaan ID %s miliknya sendiri", alumniID, id)
	} else if role == "admin" {
		// Admin can restore any
		log.Printf("Admin %s merestorasi pekerjaan ID %s", username, id)
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Anda tidak memiliki akses untuk merestorasi pekerjaan",
			"success": false,
		})
	}

	err = repository.RestorePekerjaan(db, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Pekerjaan tidak ditemukan atau tidak terhapus",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal merestorasi pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan berhasil direstorasi",
		"success": true,
	})
}

// HardDeletePekerjaanService godoc
// @Summary Hard delete riwayat pekerjaan permanen
// @Description Menghapus riwayat pekerjaan secara permanen dari database
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "Pekerjaan ID (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{} "Riwayat pekerjaan dihapus permanen"
// @Failure 403 {object} map[string]interface{} "Anda tidak punya akses"
// @Failure 404 {object} map[string]interface{} "Pekerjaan tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security Bearer
// @Router /pekerjaan/{id}/hard-delete [delete]
func HardDeletePekerjaanService(c *fiber.Ctx, db *mongo.Database) error {
	alumniID, isAlumni := c.Locals("alumni_id").(string)
	role := c.Locals("role").(string)
	username := c.Locals("username").(string)
	id := c.Params("id")

	// Get the owner of this job record
	ownerAlumniID, err := repository.GetAlumniIDByPekerjaanID(db, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Pekerjaan tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	if isAlumni && role == "alumni" {
		// Alumni trying to hard delete - can only delete their own
		if ownerAlumniID != alumniID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Maaf tidak bisa mengubah atau mengedit data selain diri sendiri",
				"success": false,
			})
		}
		log.Printf("Alumni ID %s hard delete pekerjaan ID %s miliknya sendiri", alumniID, id)
	} else if role == "admin" {
		// Admin can hard delete any
		log.Printf("Admin %s hard delete pekerjaan ID %s", username, id)
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Anda tidak memiliki akses untuk menghapus pekerjaan",
			"success": false,
		})
	}

	err = repository.HardDeletePekerjaan(db, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Pekerjaan tidak ditemukan",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan berhasil dihapus permanen",
		"success": true,
	})
}
