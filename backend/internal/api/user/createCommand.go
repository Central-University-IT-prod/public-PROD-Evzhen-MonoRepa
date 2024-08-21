package user

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CreateCommandRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ProfileID   uint   `json:"profile_id"`
}

func CreateCommand(c *fiber.Ctx) error {
	var createCommandRequest CreateCommandRequest
	if err := c.BodyParser(&createCommandRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "unable to parse body",
		})
	}

	var profile entities.Profile
	err := database.DB.Where("id = ?", createCommandRequest.ProfileID).First(&profile).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "profile with this user login and contest id not found",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	if profile.CommandID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "this profile already has command",
		})
	}

	var contest entities.Contest
	if err = database.DB.Where("id = ?", profile.ContestID).First(&contest).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	if contest.MinTeammates == 0 || contest.MaxTeammates == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "limits not set",
		})
	}

	var command entities.Command
	err = database.DB.Where("name = ? and contest_id = ?", createCommandRequest.Name, profile.ContestID).First(&command).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  false,
			"message": "command with this name already exists",
		})

	}

	command.Name = createCommandRequest.Name
	command.Description = createCommandRequest.Description
	command.ContestID = profile.ContestID
	command.OwnerID = profile.ID

	if err = database.DB.Create(&command).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	profile.CommandID = command.ID
	profile.Role = entities.Capitan
	if err = database.DB.Save(&profile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	c.Status(fiber.StatusCreated)
	return c.JSON(fiber.Map{
		"status": true,
	})
}
