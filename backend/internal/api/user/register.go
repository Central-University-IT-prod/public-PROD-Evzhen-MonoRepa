package user

import (
	"errors"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/errorz"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"unicode/utf8"
)

type Data struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {

	var registerData Data

	if err := c.BodyParser(&registerData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "unable to parse body",
		})
	}

	if err := Validate(registerData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	var user entities.User
	err := database.DB.Where("login = ?", registerData.Login).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  false,
			"message": "user with this login already exists",
		})
	}

	user.Login = registerData.Login
	user.SetPassword(registerData.Password)

	if err = database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": true,
	})
}

func Validate(user Data) error {
	if err := ValidateLogin(user.Login); err != nil {
		return err
	}
	if err := ValidatePassword(user.Password); err != nil {
		return err
	}
	return nil
}

func ValidateLogin(login string) error {
	if utf8.RuneCountInString(login) < 4 || utf8.RuneCountInString(login) > 20 {
		return errorz.ErrorInvalidLogin
	}
	return nil
}

func ValidatePassword(password string) error {
	var upper, lower, digit bool
	for _, el := range password {
		if el > 96 && el < 123 {
			lower = true
		}
		if el > 64 && el < 91 {
			upper = true
		}
		if el > 47 && el < 58 {
			digit = true
		}
	}

	if utf8.RuneCountInString(password) < 6 || utf8.RuneCountInString(password) > 100 || !upper || !lower || !digit {
		return errorz.ErrInvalidPassword
	}
	return nil
}
