package helper

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// JSONSuccess standard response
func JSONSuccess(c *fiber.Ctx, code int, data interface{}) error {
	return c.Status(code).JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

// JSONError standard error response
func JSONError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": message,
	})
}

// GetQueryInt helper: get int query with fallback
func GetQueryInt(c *fiber.Ctx, key string, fallback int) int {
	if v := c.Query(key); v != "" {
		// don't panic, parse safely
		var x int
		_, err := fmt.Sscan(v, &x)
		if err == nil {
			return x
		}
	}
	return fallback
}
