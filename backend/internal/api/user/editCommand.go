package user

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type EditCommandRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ProfileID   uint   `json:"profile_id"`
}

func EditCommand(c *fiber.Ctx) error {
	var editCommandRequest EditCommandRequest

	if err := c.BodyParser(&editCommandRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	var profile entities.Profile
	err := database.DB.Where("id = ?", editCommandRequest.ProfileID).First(&profile).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "profile with this id not found",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	if profile.Role != entities.Capitan {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  false,
			"message": "you aren't the owner",
		})
	}

	var command entities.Command
	err = database.DB.Where("id = ?", profile.CommandID).First(&command).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "command with this id not found",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	if editCommandRequest.Name != "" {
		command.Name = editCommandRequest.Name
	}
	if editCommandRequest.Description != "" {
		command.Description = editCommandRequest.Description
	}

	if err = database.DB.Save(&command).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
	})
}
