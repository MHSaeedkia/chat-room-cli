package http

import (
	"chat-room-cli/internal/http/handler"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func Run(port string) error {
	router := fiber.New()
	api := router.Group("/api/v1")

	api.Use(limiter.New(limiter.Config{
		// if ip = "127.0.0.1" , rate limmiter ignor it
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},

		Max: 10,

		Expiration: 30 * time.Second,

		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},

		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	}))

	api.Get("/health", handler.HealthCheck)

	return router.Listen(fmt.Sprintf(":%s", port))
}
