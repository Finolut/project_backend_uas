package middleware

import (
	"github.com/gofiber/fiber/v2"
	helmet "github.com/gofiber/helmet/v2"
)

func Helmet() fiber.Handler {
	return helmet.New(helmet.Config{
		// default options are sensible; customize if needed
		// e.g. ContentSecurityPolicy, CrossOriginResourcePolicy, etc.
	})
}
