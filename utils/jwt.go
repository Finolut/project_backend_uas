package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Errors exposed by this package
var (
	ErrTokenExpired   = errors.New("token expired")
	ErrTokenInvalid   = errors.New("invalid token")
	ErrTokenBadMethod = errors.New("unexpected signing method")
)

// Locals keys for fiber.Ctx.Locals
const (
	LocalsUserID    = "user_id"
	LocalsRoleID    = "role_id"
	LocalsRoleName  = "role_name"
	LocalsJWTClaims = "jwt_claims"
)

// JWTClaims used by ParseAndValidateToken and middleware.
// If you already have a JWTClaims struct in your model package, you can
// replace this with that type (and remove this type here).
type JWTClaims struct {
	UserID   string `json:"user_id"`
	RoleID   string `json:"role_id,omitempty"`
	RoleName string `json:"role_name,omitempty"`
	jwt.RegisteredClaims
}

// ParseAndValidateToken verifies and returns *JWTClaims (uses jwt package).
// Returns typed claims or error.
func ParseAndValidateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// enforce HMAC signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenBadMethod
		}
		// GetJWTSecretBytes is assumed to exist elsewhere in this package.
		// It should return the HMAC secret as []byte, e.g. []byte(os.Getenv("JWT_SECRET"))
		return GetJWTSecretBytes(), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		// jwt/v5 exposes jwt.ErrTokenExpired which can be matched with errors.Is
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, err
	}
	if !token.Valid {
		return nil, ErrTokenInvalid
	}
	return claims, nil
}

func GetJWTSecretBytes() interface{} {
	panic("unimplemented")
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
// - "jwt_claims" -> whole claims object
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
			c.Locals(LocalsRoleName, claims.RoleName)
		}
		// also store whole claims for handlers that want them
		c.Locals(LocalsJWTClaims, claims)

		return c.Next()
	}
}
