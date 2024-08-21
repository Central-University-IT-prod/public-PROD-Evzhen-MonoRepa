package user

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LeaveCommand(c *fiber.Ctx) error {
	profileID := c.Query("profile_id")

	var profile entities.Profile
	err := database.DB.Where("id = ?", profileID).First(&profile).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "profile with this id doesn't exist",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	var command entities.Command
	err = database.DB.Where("id = ?", profile.CommandID).First(&command).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "command with this id not found",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	switch profile.Role {
	case entities.Capitan:
		var teammates []entities.Profile
		if err = database.DB.Where("command_id = ?", command.ID).Find(&teammates).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		for i := range teammates {
			teammates[i].CommandID = 0
			if err = database.DB.Save(&teammates[i]).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  false,
					"message": err.Error(),
				})
			}
		}

		profile.Role = entities.Participant
		profile.CommandID = 0

		if err = database.DB.Save(&profile).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err = database.DB.Delete(&command).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	case entities.Participant:
		profile.CommandID = 0

		if err = database.DB.Save(&profile).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
	})
}
