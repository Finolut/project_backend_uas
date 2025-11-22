package middleware

import (
	"strings"

	"clean-arch/utils/mongo"

	"github.com/gofiber/fiber/v2"
)

// FileAuthRequired middleware for user file upload
// Supports both admin and regular users
func FileAuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Authorization token is required",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Invalid token format",
			})
		}

		// Try to validate as user token (works for both admin and regular users)
		userClaims, err := utils.ValidateToken(tokenParts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Token is invalid or expired",
			})
		}

		c.Locals("user_id", userClaims.UserID)
		c.Locals("username", userClaims.Username)
		c.Locals("role", userClaims.Role)

		return c.Next()
	}
}
