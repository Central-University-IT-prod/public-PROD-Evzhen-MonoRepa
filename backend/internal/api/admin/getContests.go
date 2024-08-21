package admin

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/database"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/util"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func GetContests(c *fiber.Ctx) error {
	c.Response().Header.Set("Content-Type", "application/json; charset=utf-8")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	authHeader := c.GetReqHeaders()["Authorization"][0]

	id, err := util.ParseJwt(strings.Split(authHeader, " ")[1], util.RoleAdmin)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var contests []entities.Contest
	if err = database.DB.Order("start_date ASC").Where("admin_id = ?", id).Find(&contests).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(contests)
}
