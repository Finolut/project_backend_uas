package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	mdl "clean-arch-copy/app/model/postgre"
	"clean-arch-copy/utils" // keep GenerateToken in utils
)

var (
	ErrTokenExpired   = errors.New("token expired")
	ErrTokenInvalid   = errors.New("invalid token")
	ErrTokenBadMethod = errors.New("unexpected signing method")
)

// parseToken verifies and returns *mdl.JWTClaims (uses jwt package).
// This duplicates earlier ValidateToken logic but placed inside middleware package.
// Returns typed claims or error.
func ParseAndValidateToken(tokenString string) (*mdl.JWTClaims, error) {
	claims := &mdl.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// enforce HMAC signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenBadMethod
		}
		// use utils' secret getter so secrets come from env
		return utils.GetJWTSecretBytes(), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		// inspect error for expiry
		var ve *jwt.ValidationError
		if errors.As(err, &ve) && ve.Errors&jwt.ValidationErrorExpired != 0 {
			return nil, ErrTokenExpired
		}
		return nil, err
	}
	if !token.Valid {
		return nil, ErrTokenInvalid
	}
	return claims, nil
}

// Helper: extract Bearer token from Authorization header
func extractTokenFromHeader(auth string) (string, error) {
	if auth == "" {
		return "", fiber.ErrUnauthorized
	}
	parts := strings.Fields(auth)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fiber.ErrUnauthorized
	}
	return parts[1], nil
}

// NewJWTMiddleware returns fiber middleware that validates JWT and sets locals:
// - "user_id" -> user ID (string)
// - "role_id" -> role ID (string) [if present]
// - "role_name" -> role name (string) [if present]
func NewJWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		tokStr, err := extractTokenFromHeader(auth)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing or invalid Authorization header"})
		}

		claims, err := ParseAndValidateToken(tokStr)
		if err != nil {
			// return proper HTTP code for expired token
			if errors.Is(err, ErrTokenExpired) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token expired"})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		// set locals for downstream handlers
		if claims.UserID != "" {
			c.Locals(LocalsUserID, claims.UserID)
		}
		if claims.RoleID != "" {
			c.Locals(LocalsRoleID, claims.RoleID)
		}
		if claims.RoleName != "" {
			c.Locals("role_name", claims.RoleName)
		}
		// also store whole claims for handlers that want them
		c.Locals("jwt_claims", claims)

		return c.Next()
	}
}
