package admin

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
)

func Registration(c *fiber.Ctx) error {
	// testing
	newAdmin := new(entities.Admin)
	newAdmin.Login = "admin"
	newAdmin.SetPassword("admin")
	database.DB.Create(newAdmin)
	c.Status(fiber.StatusCreated)
	return c.JSON(fiber.Map{
		"admin": newAdmin,
	})
}
