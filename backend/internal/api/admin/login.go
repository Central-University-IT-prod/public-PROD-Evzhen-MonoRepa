package admin

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/errorz"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/util"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"unicode/utf8"
)

type loginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func validateLogin(login string) error {
	if utf8.RuneCountInString(login) >= 4 && utf8.RuneCountInString(login) <= 20 {
		return nil
	} else {
		return errorz.ErrorInvalidLogin
	}
}

func Login(c *fiber.Ctx) error {
	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	var (
		admin entities.Admin
	)
	adminData := new(loginData)

	if err := c.BodyParser(adminData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errorz.CannotParseJSON.Error(),
		})
	}
	if err := validateLogin(adminData.Login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	database.DB.Where("login = ?", adminData.Login).Find(&admin)
	if admin.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": errorz.AdminNotFound.Error(),
		})
	}
	if err := admin.ComparePassword(adminData.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": errorz.PasswordIncorrect.Error(),
		})
	}

	jwt, err := util.GenerateJwt(strconv.Itoa(int(admin.ID)), util.AdminSecretKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errorz.InternalServerError.Error(),
		})
	}

	c.Status(fiber.StatusOK)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": jwt,
	})
}
