package user

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type ProfileByIDResponse struct {
	Status  bool             `json:"status"`
	Profile entities.Profile `json:"profile"`
}

func ProfileByID(c *fiber.Ctx) error {
	params := c.Params("id")

	profileID, err := strconv.Atoi(params)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "invalid profile id",
		})
	}

	var profile entities.Profile
	err = database.DB.Where("id = ?", profileID).First(&profile).Error
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

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": true,
		"body": fiber.Map{
			"profile": profile,
		},
	})
}
