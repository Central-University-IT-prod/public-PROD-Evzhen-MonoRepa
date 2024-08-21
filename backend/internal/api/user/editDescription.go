package user

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
)

type EditDescriptionRequest struct {
	Description string `json:"description"`
	ProfileID   uint   `json:"profile_id"`
}

func EditDescription(c *fiber.Ctx) error {
	var editDescriptionRequest EditDescriptionRequest

	if err := c.BodyParser(&editDescriptionRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	var profile entities.Profile
	if err := database.DB.Where("id = ?", editDescriptionRequest.ProfileID).First(&profile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	profile.Description = editDescriptionRequest.Description
	if err := database.DB.Save(&profile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		},
		)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
	})
}
