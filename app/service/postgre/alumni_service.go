package service

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"

	"clean-arch/app/model/postgre"
	"clean-arch/app/repository/postgre"

	"github.com/gofiber/fiber/v2"
)

func CheckAlumniService(c *fiber.Ctx, db *sql.DB) error {
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
		if err == sql.ErrNoRows {
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

func GetAllAlumniService(c *fiber.Ctx, db *sql.DB) error {
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

func GetAlumniByIDService(c *fiber.Ctx, db *sql.DB) error {
	username := c.Locals("username").(string)
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	log.Printf("User %s mengakses GET /api/alumni/%d", username, idInt)

	alumni, err := repository.GetAlumniByID(db, idInt)
	if err != nil {
		if err == sql.ErrNoRows {
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

func CreateAlumniService(c *fiber.Ctx, db *sql.DB) error {
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

func UpdateAlumniService(c *fiber.Ctx, db *sql.DB) error {
	username := c.Locals("username").(string)
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	log.Printf("Admin %s mengupdate alumni ID %d", username, idInt)

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

	alumni, err := repository.UpdateAlumni(db, idInt, req)
	if err != nil {
		if err == sql.ErrNoRows {
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

func DeleteAlumniService(c *fiber.Ctx, db *sql.DB) error {
	username := c.Locals("username").(string)
	id := c.Params("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	log.Printf("Admin %s menghapus alumni ID %d", username, idInt)

	err = repository.DeleteAlumni(db, idInt)
	if err != nil {
		if err == sql.ErrNoRows {
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

func GetAlumniStatisticsService(c *fiber.Ctx, db *sql.DB) error {
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

func GetAllAlumniWithPaginationService(c *fiber.Ctx, db *sql.DB) error {
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

func GetTrashedAlumniService(c *fiber.Ctx, db *sql.DB) error {
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

func SoftDeleteAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	userID, ok := c.Locals("user_id").(int)
	var deletedByID int
	if ok && userID > 0 {
		deletedByID = userID
	} else {
		// If no valid user_id, pass 0 to indicate NULL should be used
		deletedByID = 0
	}

	if err := repository.SoftDeletePekerjaanByAlumniID(db, id, deletedByID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Data pekerjaan tidak ditemukan atau sudah dihapus",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal soft delete: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan alumni berhasil dipindahkan ke trash",
		"success": true,
	})
}

func RestoreAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	if err := repository.RestoreAlumni(db, id); err != nil {
		if err == sql.ErrNoRows {
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

func HardDeleteAlumniService(c *fiber.Ctx, db *sql.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	if err := repository.HardDeletePekerjaanByAlumniID(db, id); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Hapus permanen hanya untuk data pekerjaan yang ada di trash",
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal hapus permanen: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data pekerjaan alumni dihapus permanen",
		"success": true,
	})
}
