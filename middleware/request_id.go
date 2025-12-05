package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const LocalsRequestID = "request_id"

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Set("X-Request-ID", reqID)
		c.Locals(LocalsRequestID, reqID)
		return c.Next()
	}
}
