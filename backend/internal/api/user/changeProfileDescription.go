package user

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ChangeUserProfileStruct struct {
	ProfileID   uint   `json:"profile_id"`
	Description string `json:"description"`
}

func ChangeUserProfile(c *fiber.Ctx) error {
	var (
		profile               entities.Profile
		changeUserProfileData ChangeUserProfileStruct
	)
	if err := c.BodyParser(&changeUserProfileData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to parse body",
		})
	}

	err := database.DB.Where("id = ?", changeUserProfileData.ProfileID).First(&profile).Error
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

	profile.Description = changeUserProfileData.Description
	database.DB.Save(&profile)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body": fiber.Map{
			"profile": profile,
		},
	})
}
