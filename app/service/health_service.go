package service

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckpointService(c *fiber.Ctx, db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to ping the database
	if err := db.Client().Ping(ctx, nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal cek database: " + err.Error(),
			"success": false,
		})
	}

	dbName := db.Name()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "OK",
		"success":  true,
		"database": dbName,
		"expected": "alumni_db",
		"matches":  dbName == "alumni_db",
	})
}
