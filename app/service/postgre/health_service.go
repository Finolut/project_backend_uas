package service

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func CheckpointService(c *fiber.Ctx, db *sql.DB) error {
	var currentDB string
	if err := db.QueryRow("SELECT current_database()").Scan(&currentDB); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal cek database: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "OK",
		"success":  true,
		"database": currentDB,
		"expected": "alumni_db",
		"matches":  currentDB == "alumni_db",
	})
}
