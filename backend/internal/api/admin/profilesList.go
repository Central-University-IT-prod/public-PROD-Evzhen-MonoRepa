package admin

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

func GetProfiles(c *fiber.Ctx) error {
	var (
		contest  entities.Contest
		profiles []entities.Profile
	)
	params := c.Params("contest_id")

	contestID, err := strconv.Atoi(params)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "invalid contest id",
		})
	}

	admin, err := database.GetAdminByAuthorizationToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err = database.DB.Where("id = ?", contestID).First(&contest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{})
		}
	}

	if contest.AdminID != admin.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "missing access",
		})
	}

	if err = database.DB.Where("contest_id", contestID).Find(&profiles).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(profiles)
}
