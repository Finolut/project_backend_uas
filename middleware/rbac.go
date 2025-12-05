package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// PermissionChecker adalah fungsi yang mengecek apakah role memiliki permission.
// Implementasikan wrapper yang memanggil RBACService.HasPermissionByRoleID di tempat wiring.
type PermissionChecker func(roleID string, permission string) (bool, error)

// RequirePermission returns a middleware that checks permission string (e.g. "achievement:verify")
func RequirePermission(check PermissionChecker, permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals(LocalsRoleID).(string)
		if !ok || role == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "role not found in token"})
		}
		okPerm, err := check(role, permission)
		if err != nil {
			// optionally log error
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "permission check failed"})
		}
		if !okPerm {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "permission denied"})
		}
		return c.Next()
	}
}
