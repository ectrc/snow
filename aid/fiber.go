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

func FiberLimiter(n int) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        n,
		Expiration: 1 * time.Minute,
	})
}

func FiberCors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Requested-With",
	})
}

// https://github.com/gofiber/fiber/issues/510
func FiberGetQueries(c *fiber.Ctx, queryKeys ...string) map[string][]string {
	argsMaps := make(map[string][]string)
	for _, keys := range queryKeys {
		param := c.Request().URI().QueryArgs().PeekMulti(keys)
		for _, value := range param {
			argsMaps[keys] = append(argsMaps[keys], string(value))
		}
	}
	return argsMaps
}