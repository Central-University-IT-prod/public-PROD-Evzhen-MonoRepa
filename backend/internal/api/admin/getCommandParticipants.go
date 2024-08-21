package admin

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetCommandParticipants(c *fiber.Ctx) error {
	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	params := c.Params("command_id")

	commandID, err := strconv.Atoi(params)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "invalid command id",
		})
	}

	teammates := database.GetCommandParticipants(uint(commandID))
	return c.Status(fiber.StatusOK).JSON(teammates)
}
