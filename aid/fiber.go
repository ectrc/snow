package aid

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func FiberLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format: "(${method}) (${status}) (${latency}) ${path}\n",
		Next: func(c *fiber.Ctx) bool {
			return c.Response().StatusCode() == 302
		},
	})
}

func FiberLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	})
}

func FiberCors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Requested-With",
	})
}