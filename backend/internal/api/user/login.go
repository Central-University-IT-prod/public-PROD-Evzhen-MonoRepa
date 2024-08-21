package user

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/util"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

func Login(c *fiber.Ctx) error {
	var loginData Data

	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "unable to parse body",
		})
	}

	if err := Validate(loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	var user entities.User
	err := database.DB.Where("login = ?", loginData.Login).First(&user).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "user with this login and password not found",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	if err = user.ComparePassword(loginData.Password); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "user with this login and password not found",
		})
	}

	jwt, err := util.GenerateJwt(strconv.Itoa(int(user.ID)), util.RoleUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"body": fiber.Map{
			"token": jwt,
		},
	})
}
