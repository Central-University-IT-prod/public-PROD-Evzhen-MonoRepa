package routes

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/routes/admin"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/routes/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Setup(app *fiber.App) {
	app.Use(cors.New())
	api := app.Group("/api", cors.New())
	admin.Setup(api)
	user.Setup(api)
	app.Get("/ping", func(c *fiber.Ctx) error {
		c.Status(fiber.StatusOK)
		return c.JSON("pong")
	})
}
