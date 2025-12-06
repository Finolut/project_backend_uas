package middleware

import (
	"strings"
	"time"

	"clean-arch-copy/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Key names for locals
const (
	LocalsUserID = "user_id"
	LocalsRoleID = "role_id"
)

// NewJWTMiddleware returns a Fiber middleware that validates JWT and sets c.Locals("user_id", id)
func NewJWTMiddleware() fiber.Handler {
	secret := config.Get().JWTSecret
	if secret == "" {
		secret = "dev-secret"
	}
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
		}
		parts := strings.Fields(auth)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header"})
		}
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// ensure signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secret), nil
		}, jwt.WithLeeway(5*time.Second))
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}
		// expected claims: sub (user id), role
		if sub, ok := claims["sub"].(string); ok && sub != "" {
			c.Locals(LocalsUserID, sub)
		}
		if role, ok := claims["role"].(string); ok && role != "" {
			c.Locals(LocalsRoleID, role)
		}
		return c.Next()
	}
}
