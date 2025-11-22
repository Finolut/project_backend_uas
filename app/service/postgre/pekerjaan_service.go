package service

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"clean-arch/app/model/postgre"
	"clean-arch/app/repository/postgre"

	"github.com/gofiber/fiber/v2"
)

func GetAllPekerjaanService(c *fiber.Ctx, db *sql.DB) error {
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

func GetPekerjaanByIDService(c *fiber.Ctx, db *sql.DB) error {
	username := c.Locals("username").(string)
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	log.Printf("User %s mengakses GET /api/pekerjaan/%d", username, idInt)

	pekerjaan, err := repository.GetPekerjaanByID(db, idInt)
	if err != nil {
		if err == sql.ErrNoRows {
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

func GetPekerjaanByAlumniIDService(c *fiber.Ctx, db *sql.DB) error {
	username := c.Locals("username").(string)
	alumniID := c.Params("alumni_id")

	alumniIDInt, err := strconv.Atoi(alumniID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Alumni ID tidak valid",
			"success": false,
		})
	}

	log.Printf("Admin %s mengakses GET /api/pekerjaan/alumni/%d", username, alumniIDInt)

	pekerjaan, err := repository.GetPekerjaanByAlumniID(db, alumniIDInt)
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

func CreatePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	alumniID := c.Locals("alumni_id").(int)
	nama := c.Locals("nama").(string)
	log.Printf("Alumni %s (ID: %d) menambah pekerjaan baru", nama, alumniID)

	var req model.CreatePekerjaanRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format data tidak valid: " + err.Error(),
			"success": false,
		})
	}

	req.AlumniID = alumniID

	// Validasi input
	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Nama perusahaan dan posisi jabatan wajib diisi",
			"success": false,
		})
	}

	pekerjaan, err := repository.CreatePekerjaan(db, req)
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

func UpdatePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	alumniID := c.Locals("alumni_id").(int)
	nama := c.Locals("nama").(string)
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	log.Printf("Alumni %s (ID: %d) mengupdate pekerjaan ID %d", nama, alumniID, idInt)

	ownerAlumniID, err := repository.GetAlumniIDByPekerjaanID(db, idInt)
	if err != nil {
		if err == sql.ErrNoRows {
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

	pekerjaan, err := repository.UpdatePekerjaan(db, idInt, req)
	if err != nil {
		if err == sql.ErrNoRows {
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

func DeletePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	username := c.Locals("username").(string)
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	log.Printf("Admin %s menghapus pekerjaan ID %d", username, idInt)

	err = repository.DeletePekerjaan(db, idInt)
	if err != nil {
		if err == sql.ErrNoRows {
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

func GetAllPekerjaanWithPaginationService(c *fiber.Ctx, db *sql.DB) error {
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

func SoftDeletePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	alumniID, isAlumni := c.Locals("alumni_id").(int)
	role := c.Locals("role").(string)
	username := c.Locals("username").(string)
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	// Get the owner of this job record
	ownerAlumniID, err := repository.GetAlumniIDByPekerjaanID(db, idInt)
	if err != nil {
		if err == sql.ErrNoRows {
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
		log.Printf("Alumni ID %d menghapus pekerjaan ID %d miliknya sendiri", alumniID, idInt)
	} else if role == "admin" {
		// Admin can delete any
		log.Printf("Admin %s menghapus pekerjaan ID %d", username, idInt)
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Anda tidak memiliki akses untuk menghapus pekerjaan",
			"success": false,
		})
	}

	var deletedBy int
	if role == "admin" {
		// For admin, we need to get the user_id from the admin token
		userID, ok := c.Locals("user_id").(int)
		if !ok {
			userID = -1
		}
		deletedBy = userID
	} else {
		deletedBy = alumniID
	}

	err = repository.SoftDeletePekerjaan(db, idInt, deletedBy)
	if err != nil {
		if err == sql.ErrNoRows {
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

func RestorePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	alumniID, isAlumni := c.Locals("alumni_id").(int)
	role := c.Locals("role").(string)
	username := c.Locals("username").(string)
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	// Get the owner of this job record
	ownerAlumniID, err := repository.GetAlumniIDByPekerjaanID(db, idInt)
	if err != nil {
		if err == sql.ErrNoRows {
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
		log.Printf("Alumni ID %d merestorasi pekerjaan ID %d miliknya sendiri", alumniID, idInt)
	} else if role == "admin" {
		// Admin can restore any
		log.Printf("Admin %s merestorasi pekerjaan ID %d", username, idInt)
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Anda tidak memiliki akses untuk merestorasi pekerjaan",
			"success": false,
		})
	}

	err = repository.RestorePekerjaan(db, idInt)
	if err != nil {
		if err == sql.ErrNoRows {
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

func HardDeletePekerjaanService(c *fiber.Ctx, db *sql.DB) error {
	alumniID, isAlumni := c.Locals("alumni_id").(int)
	role := c.Locals("role").(string)
	username := c.Locals("username").(string)
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	// Get the owner of this job record
	ownerAlumniID, err := repository.GetAlumniIDByPekerjaanID(db, idInt)
	if err != nil {
		if err == sql.ErrNoRows {
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
		log.Printf("Alumni ID %d hard delete pekerjaan ID %d miliknya sendiri", alumniID, idInt)
	} else if role == "admin" {
		// Admin can hard delete any
		log.Printf("Admin %s hard delete pekerjaan ID %d", username, idInt)
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Anda tidak memiliki akses untuk menghapus pekerjaan",
			"success": false,
		})
	}

	err = repository.HardDeletePekerjaan(db, idInt)
	if err != nil {
		if err == sql.ErrNoRows {
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
