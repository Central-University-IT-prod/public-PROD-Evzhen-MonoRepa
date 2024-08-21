package admin

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ApproveCommand(c *fiber.Ctx) error {
	var (
		command    entities.Command
		contest    entities.Contest
		isApproved bool
	)
	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	commandID := c.Params("command_id")
	approved := c.Params("approved")

	switch approved {
	case "1":
		isApproved = true
	case "0":
		isApproved = false
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid approved data",
		})
	}

	admin, err := database.GetAdminByAuthorizationToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.DB.Where("id = ?", commandID).First(&command).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "command not found",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.DB.Where("id = ?", command.ContestID).First(&contest).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "contest not found",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if admin.ID != contest.AdminID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "missing access",
		})
	}

	command.Approved = isApproved
	database.DB.Save(&command)
	return c.Status(fiber.StatusOK).JSON(command)
}
