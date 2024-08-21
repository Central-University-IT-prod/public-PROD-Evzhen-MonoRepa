package admin

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

func DeleteProfile(c *fiber.Ctx) error {
	var (
		contest entities.Contest
		profile entities.Profile
	)

	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	params := c.Params("id")

	profileID, err := strconv.Atoi(params)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "invalid profile id",
		})
	}

	admin, err := database.GetAdminByAuthorizationToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err = database.DB.Where("id = ?", profileID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if err = database.DB.Where("id = ?", profile.ContestID).First(&contest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if contest.AdminID != admin.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "missing access",
		})
	}

	if err = database.DB.Delete(&profile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//TODO: delete command if owner

	return c.SendStatus(fiber.StatusOK)
}
