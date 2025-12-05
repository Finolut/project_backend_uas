package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// NewFiberApp returns a configured *fiber.App with standard middleware.
// Optionally pass a custom logger writer (from InitLogger) to wire into fiber logger.
func NewFiberApp(logWriter ...interface{}) *fiber.App {
	app := fiber.New(fiber.Config{
		// you can set ReadTimeout/WriteTimeout here if you want
		// ReadTimeout:  5 * time.Second,
		// WriteTimeout: 10 * time.Second,
	})

	// Recover from panics
	app.Use(recover.New())

	// CORS: adjust as needed
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Logger middleware: if logWriter provided and is io.Writer, use it
	if len(logWriter) > 0 {
		// fiber's logger middleware accepts config with Output
		if w, ok := logWriter[0].(interface{ Write([]byte) (int, error) }); ok {
			app.Use(fiberlogger.New(fiberlogger.Config{
				TimeFormat: "02-Jan-2006 15:04:05",
				Format:     "${time} | ${status} | ${method} ${path} - ${ip} - ${latency}\n",
				Output:     w,
				TimeZone:   "Local",
			}))
		} else {
			// fallback to default logger
			app.Use(fiberlogger.New())
		}
	} else {
		app.Use(fiberlogger.New())
	}

	// Optional: add security headers, rate limiter, helmet-like middleware, etc.
	// app.Use(helmet.New()) // if you add helmet or similar package

	// Example: set global timeout on handlers (if desired)
	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Server", "clean-arch-app")
		// you may set deadline here by wrapping context but Fiber handlers use its own ctx
		return c.Next()
	})

	return app
}
